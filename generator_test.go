package parcello_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/parcello"
	"github.com/phogolabs/parcello/fake"
)

var _ = Describe("Generator", func() {
	var (
		generator  *parcello.Generator
		bundle     *parcello.Bundle
		node       *parcello.Node
		buffer     *parcello.ResourceFile
		fileSystem *fake.FileSystem
	)

	BeforeEach(func() {
		bundle = &parcello.Bundle{
			Name: "bundle",
			Body: []byte{
				31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 212, 146, 223, 171, 218,
				48, 28, 197, 251, 156, 191, 226, 60, 110, 96, 103, 90, 107,
				156, 250, 212, 185, 12, 202, 172, 150, 54, 3, 125, 146, 216,
				70, 45, 104, 219, 37, 41, 251, 247, 71, 117, 212, 135, 93, 132,
				43, 23, 46, 247, 243, 114, 242, 227, 36, 57, 95, 242, 45, 164,
				149, 123, 105, 212, 240, 82, 30, 181, 180, 101, 93, 13, 41,
				165, 140, 142, 253, 128, 94, 217, 25, 101, 219, 230, 139, 249,
				125, 118, 158, 165, 187, 134, 5, 193, 85, 39, 108, 124, 85,
				234, 223, 230, 29, 1, 27, 59, 222, 200, 103, 62, 155, 140, 188,
				17, 115, 168, 79, 3, 202, 28, 208, 167, 95, 124, 5, 173, 177,
				82, 59, 148, 106, 121, 206, 79, 15, 124, 198, 202, 195, 225,
				193, 254, 191, 90, 122, 253, 32, 184, 46, 194, 214, 214, 238,
				81, 85, 74, 75, 171, 10, 72, 139, 184, 174, 16, 54, 26, 152,
				194, 99, 51, 111, 58, 243, 40, 22, 60, 19, 240, 169, 247, 149,
				184, 46, 146, 179, 146, 70, 161, 168, 81, 213, 22, 249, 73,
				86, 71, 5, 123, 82, 168, 228, 69, 65, 90, 171, 203, 125, 107,
				149, 33, 157, 185, 91, 155, 161, 109, 8, 89, 164, 60, 20, 28,
				34, 252, 182, 228, 136, 126, 96, 181, 22, 224, 155, 40, 19,
				25, 250, 246, 51, 248, 68, 80, 22, 232, 17, 124, 35, 110, 163,
				206, 190, 250, 181, 92, 34, 73, 163, 56, 76, 183, 248, 201,
				183, 3, 130, 66, 153, 92, 151, 77, 119, 248, 5, 243, 128, 32,
				215, 170, 171, 108, 39, 45, 32, 162, 152, 103, 34, 140, 147,
				222, 64, 62, 207, 239, 41, 139, 250, 79, 69, 200, 247, 116,
				157, 220, 83, 254, 151, 112, 254, 222, 159, 246, 134, 252, 5,
				0, 0, 255, 255, 194, 146, 255, 65, 145, 12, 202, 128, 134, 150,
				6, 6, 134, 166, 241, 165, 197, 169, 69, 197, 52, 205, 255, 198,
				230, 38, 240, 252, 111, 96, 110, 12, 206, 255, 198, 70, 163,
				249, 159, 30, 0, 107, 254, 119, 43, 202, 132, 228, 127, 51,
				5, 67, 75, 43, 3, 3, 43, 67, 83, 106, 228, 127, 148, 236, 15,
				78, 86, 160, 188, 14, 202, 236, 158, 126, 33, 200, 153, 26,
				57, 243, 42, 164, 101, 22, 21, 151, 196, 131, 13, 6, 231, 110,
				100, 185, 156, 68, 100, 41, 80, 78, 70, 203, 202, 88, 115, 50,
				216, 106, 107, 46, 174, 129, 14, 251, 193, 0, 0, 0, 0, 0, 255,
				255, 130, 231, 127, 72, 25, 138, 148, 249, 13, 13, 141, 205,
				41, 202, 246, 112, 64, 40, 255, 27, 26, 155, 162, 230, 127,
				67, 51, 83, 35, 211, 209, 252, 79, 15, 64, 68, 254, 55, 52,
				180, 50, 54, 71, 205, 255, 144, 28, 86, 156, 145, 95, 174, 11,
				206, 76, 92, 193, 174, 62, 174, 206, 33, 10, 90, 10, 110, 65,
				254, 190, 163, 25, 108, 232, 0, 0, 0, 0, 0, 255, 255,
			},
		}

		node = &parcello.Node{
			Name:    "resource",
			Content: &[]byte{},
			Mutex:   &sync.RWMutex{},
		}

		buffer = parcello.NewResourceFile(node)

		fileSystem = &fake.FileSystem{}
		fileSystem.OpenFileReturns(buffer, nil)

		generator = &parcello.Generator{
			FileSystem: fileSystem,
			Config: &parcello.GeneratorConfig{
				Package: "mypackage",
			},
		}
	})

	It("writes the bundle to the destination successfully", func() {
		Expect(generator.Compose(bundle)).To(Succeed())
		Expect(fileSystem.OpenFileCallCount()).To(Equal(1))

		filename, flag, mode := fileSystem.OpenFileArgsForCall(0)
		Expect(filename).To(Equal("bundle.go"))
		Expect(flag).To(Equal(os.O_WRONLY | os.O_CREATE | os.O_TRUNC))
		Expect(mode).To(Equal(os.FileMode(0600)))

		_, err := buffer.Seek(0, io.SeekStart)
		Expect(err).To(BeNil())
		content, err := ioutil.ReadAll(buffer)
		Expect(err).To(BeNil())

		Expect(content).To(ContainSubstring("package mypackage"))
		Expect(content).To(ContainSubstring("func init()"))
		Expect(content).To(ContainSubstring("parcello.AddResource"))
		Expect(content).NotTo(ContainSubstring("// Code generated by parcello; DO NOT EDIT."))
	})

	Context("when include API documentation is enabled", func() {
		BeforeEach(func() {
			generator.Config.InlcudeDocs = true
		})

		It("includes the documentation", func() {
			Expect(generator.Compose(bundle)).To(Succeed())

			_, err := buffer.Seek(0, io.SeekStart)
			Expect(err).To(BeNil())
			content, err := ioutil.ReadAll(buffer)
			Expect(err).To(BeNil())

			Expect(content).To(ContainSubstring("package mypackage"))
			Expect(content).To(ContainSubstring("func init()"))
			Expect(content).To(ContainSubstring("parcello.AddResource"))
			Expect(content).To(ContainSubstring("// Code generated by parcello; DO NOT EDIT."))
		})
	})

	Context("when the package name is not provided", func() {
		BeforeEach(func() {
			generator.Config.Package = ""
		})

		It("returns the error", func() {
			Expect(generator.Compose(bundle)).To(MatchError("3:1: expected 'IDENT', found 'import'"))
		})
	})

	Context("when the file system fails", func() {
		It("returns the error", func() {
			fileSystem.OpenFileReturns(nil, fmt.Errorf("Oh no!"))
			Expect(generator.Compose(bundle)).To(MatchError("Oh no!"))
		})
	})

	Context("when writing the bundle fails", func() {
		It("returns the error", func() {
			buffer := &fake.File{}
			buffer.WriteReturns(0, fmt.Errorf("Oh no!"))
			fileSystem.OpenFileReturns(buffer, nil)

			Expect(generator.Compose(bundle)).To(MatchError("Oh no!"))
		})
	})
})
