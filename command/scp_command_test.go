package command

import (
	"10gen.com/mci"
	"10gen.com/mci/util"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestScpCommand(t *testing.T) {

	Convey("With files to scp", t, func() {

		// the local files and target directory for scping
		mciHome, err := mci.FindMCIHome()
		So(err, ShouldBeNil)

		tmpBase := filepath.Join(mciHome, "command/testdata/tmp")
		fileToScp := filepath.Join(tmpBase, "copy_me_please.txt")
		directoryToScp := filepath.Join(tmpBase, "copy_my_children_please")
		nestedFileToScp := filepath.Join(directoryToScp,
			"copy_me_too_please.txt")
		targetDirectory := filepath.Join(tmpBase, "feed_me_files")

		// remove the files and directories, if they exist (start clean)
		exists, err := util.FileExists(tmpBase)
		So(err, ShouldBeNil)
		if exists {
			So(os.RemoveAll(tmpBase), ShouldBeNil)
		}
		So(os.MkdirAll(tmpBase, 0777), ShouldBeNil)

		// prevent permissions issues
		syscall.Umask(0000)

		// create the files / directories to be used
		So(ioutil.WriteFile(fileToScp, []byte("hello"), 0777), ShouldBeNil)
		So(os.Mkdir(directoryToScp, 0777), ShouldBeNil)
		So(ioutil.WriteFile(nestedFileToScp, []byte("hi"), 0777), ShouldBeNil)
		So(os.Mkdir(targetDirectory, 0777), ShouldBeNil)

		Convey("when running scp commands", func() {

			Convey("copying files should work in both directions (local to"+
				" remote and remote to local)", func() {

				// scp the file from local to remote
				scpCommand := &ScpCommand{
					Source:         fileToScp,
					Dest:           targetDirectory,
					Stdout:         ioutil.Discard,
					Stderr:         ioutil.Discard,
					RemoteHostName: TestRemote,
					User:           TestRemoteUser,
					Options:        []string{"-i", TestRemoteKey},
					SourceIsRemote: false,
				}
				So(scpCommand.Run(), ShouldBeNil)

				// make sure the file was scp-ed over
				newFileContents, err := ioutil.ReadFile(
					filepath.Join(targetDirectory, "copy_me_please.txt"))
				So(err, ShouldBeNil)
				So(newFileContents, ShouldResemble, []byte("hello"))

				// remove the file
				So(os.Remove(filepath.Join(targetDirectory,
					"copy_me_please.txt")), ShouldBeNil)

				// scp the file from remote to local
				scpCommand = &ScpCommand{
					Source:         fileToScp,
					Dest:           targetDirectory,
					Stdout:         ioutil.Discard,
					Stderr:         ioutil.Discard,
					RemoteHostName: TestRemote,
					User:           TestRemoteUser,
					Options:        []string{"-i", TestRemoteKey},
					SourceIsRemote: true,
				}
				So(scpCommand.Run(), ShouldBeNil)

				// make sure the file was scp-ed over
				newFileContents, err = ioutil.ReadFile(
					filepath.Join(targetDirectory, "copy_me_please.txt"))
				So(err, ShouldBeNil)
				So(newFileContents, ShouldResemble, []byte("hello"))

			})

			Convey("additional scp options should be passed correctly to the"+
				" command", func() {

				// scp recursively, using the -r flag
				scpCommand := &ScpCommand{
					Source:         directoryToScp,
					Dest:           targetDirectory,
					Stdout:         ioutil.Discard,
					Stderr:         ioutil.Discard,
					RemoteHostName: TestRemote,
					User:           TestRemoteUser,
					Options:        []string{"-i", TestRemoteKey, "-r"},
					SourceIsRemote: false,
				}
				So(scpCommand.Run(), ShouldBeNil)

				// make sure the entire directory was scp-ed over
				nestedFileContents, err := ioutil.ReadFile(
					filepath.Join(targetDirectory, "copy_my_children_please",
						"copy_me_too_please.txt"))
				So(err, ShouldBeNil)
				So(nestedFileContents, ShouldResemble, []byte("hi"))
			})

		})

	})
}