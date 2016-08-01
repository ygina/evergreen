package model

import (
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// helper for returning the variant representing a matrix value.
func findMatrixVariant(bvs []BuildVariant, cell matrixValue) *BuildVariant {
	for i, v := range bvs {
		found := 0
		for key, val := range cell {
			if x, ok := v.Expansions[key]; ok && x == val {
				found++
			}
		}
		if found == len(cell) {
			return &bvs[i]
		}
	}
	return nil
}

// helper for turning a list of a variant's tasks into a list of names.
func taskNames(v *BuildVariant) []string {
	var names []string
	for _, t := range v.Tasks {
		names = append(names, t.Name)
	}
	return names
}

func TestPythonMatrixIntegration(t *testing.T) {
	Convey("With a sample matrix project mocking up a python driver", t, func() {
		p := Project{}
		bytes, err := ioutil.ReadFile("testdata/matrix_python.yml")
		So(err, ShouldBeNil)
		Convey("the project should parse properly", func() {
			err := LoadProjectInto(bytes, "python", &p)
			So(err, ShouldBeNil)
			Convey("and contain the correct variants", func() {
				So(len(p.BuildVariants), ShouldEqual, (2*2*4 - 4))
				Convey("so that excluded matrix cells are not created", func() {
					So(findMatrixVariant(p.BuildVariants, matrixValue{
						"os": "windows", "python": "pypy", "c-extensions": "with-c",
					}), ShouldBeNil)
					So(findMatrixVariant(p.BuildVariants, matrixValue{
						"os": "windows", "python": "jython", "c-extensions": "with-c",
					}), ShouldBeNil)
					So(findMatrixVariant(p.BuildVariants, matrixValue{
						"os": "linux", "python": "pypy", "c-extensions": "with-c",
					}), ShouldBeNil)
					So(findMatrixVariant(p.BuildVariants, matrixValue{
						"os": "linux", "python": "jython", "c-extensions": "with-c",
					}), ShouldBeNil)
				})
				Convey("so that Windows builds without C extensions exclude LDAP tasks", func() {
					v := findMatrixVariant(p.BuildVariants, matrixValue{
						"os":           "windows",
						"python":       "python3",
						"c-extensions": "without-c",
					})
					So(v, ShouldNotBeNil)
					tasks := taskNames(v)
					So(len(tasks), ShouldEqual, 7)
					So(tasks, ShouldNotContain, "ldap_auth")
				})
				Convey("so that the linux/python3/c variant has a lint task", func() {
					v := findMatrixVariant(p.BuildVariants, matrixValue{
						"os":           "linux",
						"python":       "python3",
						"c-extensions": "with-c",
					})
					So(v, ShouldNotBeNil)
					tasks := taskNames(v)
					So(len(tasks), ShouldEqual, 9)
					So(tasks, ShouldContain, "ldap_auth")
					So(tasks, ShouldContain, "lint")
				})
			})
		})
	})
}
