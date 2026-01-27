package client_test

// You MUST NOT change these default imports.  ANY additional imports may
// break the autograder and everyone will be sad.

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports.
	_ "encoding/hex"
	_ "errors"
	_ "strconv"
	_ "strings"
	"testing"

	"github.com/google/uuid"
	_ "github.com/google/uuid"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect().
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	userlib "github.com/cs161-staff/project2-userlib"

	"github.com/cs161-staff/project2-starter-code/client"
)

func TestSetupAndExecution(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Tests")
}

// ================================================
// Global Variables (feel free to add more!)
// ================================================
const defaultPassword = "password"
const emptyString = ""
const contentOne = "Bitcoin is Nick's favorite "
const contentTwo = "digital "
const contentThree = "cryptocurrency!"

// ================================================
// Describe(...) blocks help you organize your tests
// into functional categories. They can be nested into
// a tree-like structure.
// ================================================

var _ = Describe("Client Tests", func() {

	// A few user declarations that may be used for testing. Remember to initialize these before you
	// attempt to use them!
	var alice *client.User
	var bob *client.User
	var charles *client.User
	// var doris *client.User
	// var eve *client.User
	// var frank *client.User
	// var grace *client.User
	// var horace *client.User
	// var ira *client.User

	// These declarations may be useful for multi-session testing.
	var alicePhone *client.User
	var aliceLaptop *client.User
	var aliceDesktop *client.User

	var err error

	// A bunch of filenames that may be useful.
	aliceFile := "aliceFile.txt"
	bobFile := "bobFile.txt"
	charlesFile := "charlesFile.txt"
	// dorisFile := "dorisFile.txt"
	// eveFile := "eveFile.txt"
	// frankFile := "frankFile.txt"
	// graceFile := "graceFile.txt"
	// horaceFile := "horaceFile.txt"
	// iraFile := "iraFile.txt"

	BeforeEach(func() {
		// This runs before each test within this Describe block (including nested tests).
		// Here, we reset the state of Datastore and Keystore so that tests do not interfere with each other.
		// We also initialize
		userlib.DatastoreClear()
		userlib.KeystoreClear()
	})

	Describe("Basic Tests", func() {

		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})

		Specify("Basic Test: Testing Single User Store/Load/Append.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Create/Accept Invite Functionality with multiple users and multiple instances.", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			aliceDesktop, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop storing file %s with content: %s", aliceFile, contentOne)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for Bob.")
			invite, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob appending to file %s, content: %s", bobFile, contentTwo)
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop appending to file %s, content: %s", aliceFile, contentThree)
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that aliceLaptop sees expected file data.")
			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Bob sees expected file data.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Getting third instance of Alice - alicePhone.")
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that alicePhone sees Alice's changes.")
			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Revoke Functionality", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Bob creating invite for Charles for file %s, and Charlie accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Charles can load the file.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			_, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that the revoked users cannot append to the file.")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			err = charles.AppendToFile(charlesFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

	})

	// //------------------------------------------------Caleb's tests-----------------------------------------------------------------------------
	// // 1. Tests Initialize and Gets on User other than Alice
	// // 2. Checks to see if Empty Username errors out
	// // 3. Checks to see if duplicate Username Errors out
	// // 4. Checks to see if Wrong Passwoerd Errors Out
	// // 5. Tests to see if GetUser errors out on Non-existent user
	// // 6. Checks if File gets overwritten properly
	// // 7. Load File Should Fail on non existent files
	// // 8. Append File Should Fail on non existent files
	// // 9. StoreFile should be able to overwrite a file after append
	// // 10. Storing and Loading should work with Empty Files
	// // 11. Appending should work with Empty Files
	// // 12. Tests if createinvitation errors on non existent files as it should
	// // 13. Testing to see if create invitation errors on non existent users
	// // 14. testing if accept invitation fails on existing filensames
	// // 15. testing if someone who is not the owener (invited) can not revoke other
	// // 16. testing to make sure someone with no access (someone who has been revoked) can not create new invitationGI
	// // 17. Checks if multi devices work (like overwriting reflects on different devices)
	// // 18. Checking loadfile doesnt work for tamepred files
	// // 19. Checking acceptinivtation doesnt work if the invitation been tamepred with

	Describe("Our Own Tests", func() { // Can Also Check if it properly stops Alice from creating another Alice user like no double
		// 1. Tests Initialize and Gets on User other than Alice
		Specify("Testing InitUser/GetUser on other Users besides Alice.", func() {
			// Creating (Initializing User Bob)
			userlib.DebugMsg("Initializing User Other Than Alice")
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			// Checks if Bob was created successfully if not it will ERror Out
			userlib.DebugMsg("Getting user Bob.")
			bob, err = client.GetUser("bob", defaultPassword)

			// Check login was successful (No error returned)
			Expect(err).To(BeNil())
		})

		// 2. Checks to see if Empty Username errors out
		// Specify("Username is an empty username when Initializing.", func() {
		// 	// Creates new User using InitUser, but the user is empty
		// 	_, err = client.InitUser(emptyString, defaultPassword)

		// 	// Expects an Error to happen (As blank username should not work)
		// 	Expect(err).ToNot(BeNil())
		// })

		Specify("Username is an empty username when Initializing.", func() {
			// Creates new User using InitUser, but the user is empty
			_, err = client.InitUser(emptyString, defaultPassword)

			// Expects an Error to happen (As blank username should work)
			Expect(err).ToNot(BeNil())
		})

		// 3. Checks to see if duplicate Username Errors out
		Specify("Testing if it properly Errors out When using InitUser on same username.", func() {
			userlib.DebugMsg("Initializing User Alice")
			// Create user Alice
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			Expect(alice).ToNot(BeNil())

			// Creates user alice again
			userlib.DebugMsg("Initializing Alice (Duplicate User) Again.")
			var duplicateAlice *client.User

			// Tries to create another user alice with a different password (Password does not matter really ngl)
			duplicateAlice, err = client.InitUser("alice", "another_password")

			// Expect test to error out because should not be able to create same username
			Expect(err).ToNot(BeNil())

			// Make sure duplicate user was not created
			Expect(duplicateAlice).To(BeNil())
		})

		// 4. Checks to see if Wrong Passwoerd Errors Out
		Specify("Wrong Password.", func() {
			userlib.DebugMsg("Initializing user Alice")

			// Creates new user Alice
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			userlib.DebugMsg("Attempt to Login to Alice with Incorrect Password.")

			// New Variable to store result of GetUser
			var aliceWrongPass *client.User

			// Tries getting the user alice with the wrong password
			aliceWrongPass, err = client.GetUser("alice", "wrong-password")

			// Expecting the test to error out as it should not be possible
			Expect(err).ToNot(BeNil())
			Expect(aliceWrongPass).To(BeNil())
		})

		// 5. Tests to see if GetUser errors out on Non-existent user
		Specify("Testing GetUser on Non-Existent Users.", func() {
			userlib.DebugMsg("Attempt to Initialize GetUser on a Non-Existent user.")

			// Variable to store result of GetUesr
			var nonExistentUser *client.User

			// Tries using GetUser on a non-existant user (shouldnt not be possible)
			nonExistentUser, err = client.GetUser("non_existent", "random_password")
			Expect(err).ToNot(BeNil())
			Expect(nonExistentUser).To(BeNil())
		})

		// Befor eEach(func() {
		// 	alice, err = client.InitUser("alice", defaultPassword)
		// 	Expect(err).To(BeNil())
		// })
	})

	Describe("Load, Store, Append", func() {
		BeforeEach(func() {
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})

		// 6. Checks if File gets overwritten properly
		Specify("Testing to see if Store File correctly overwrites existing File.", func() {
			// First Store Content
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			// Should be new content
			userlib.DebugMsg("Overwriting Existing File")
			err = alice.StoreFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			// Loading and Checking if File got overwritten
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentThree))) // Should be the new content which is contentThree
		})

		// 7. Load File Should Fail on non existent files
		Specify("Testing to see if LoadFile successfully fails on non existent files.", func() {
			_, err = alice.LoadFile("non_existent.txt")
			Expect(err).ToNot(BeNil())
		})

		// 8. Append File Should Fail on non existent files
		Specify("Testing to see if Append To File successfully fails on non existent files.", func() {
			err = alice.AppendToFile("non_existent.txt", []byte(contentOne))
			Expect(err).ToNot(BeNil())
		})

		// 9. StoreFile should be able to overwrite file after append
		Specify("Testing to see if StoreFile overwriting works correctly after appending.", func() {
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			// Overwriting the appeneded file
			userlib.DebugMsg("Overwrite Appended File.")
			err = alice.StoreFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			// Checking
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentThree)))
		})

		// 10. Storing and Loading should work with Empty Files
		Specify("Testing to see if storing and loading works with empty files.", func() {
			err = alice.StoreFile(aliceFile, []byte(emptyString))
			Expect(err).To(BeNil())
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(emptyString)))
		})

		// 11. Appending should work with Empty Files
		Specify("TEsting to see if appending works with empty files.", func() {
			err = alice.StoreFile(aliceFile, []byte(emptyString))
			Expect(err).To(BeNil())
			err = alice.AppendToFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))
		})

		// BeforeEach(func() {
		// 	alice, err = client.InitUser("alice", defaultPassword)
		// 	Expect(err).To(BeNil())

		// 	bob, err = client.InitUser("bob", defaultPassword)
		// 	Expect(err).To(BeNil())

		// 	err = alice.StoreFile(aliceFile, []byte(contentOne))
		// 	Expect(err).To(BeNil())
		// })
	})

	Describe("CreateInv and Accept Inv", func() {
		BeforeEach(func() {
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())
		})

		// 12. Tests if createinvitation errors on non existent files as it should
		Specify("Testing CreateInvitation failing on non existent files.", func() {
			_, err = alice.CreateInvitation("non_existent.txt", "bob")
			Expect(err).ToNot(BeNil())
		})

		// 13. Testing to see if create invitation errors on non existent users
		Specify("Testing to see if CreateInivitation fails on non existent receipients.", func() {
			_, err = alice.CreateInvitation(aliceFile, "non_existent_user")
			Expect(err).ToNot(BeNil())
		})

		// 14. testing if accept invitation fails on existing filensames
		Specify("TEsting to see if AcceptInivitation fails if filename already exists.", func() {
			err = bob.StoreFile(bobFile, []byte("bobs_stuff"))
			Expect(err).To(BeNil())

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			// bob tries accepting an invitation of a filename he already has should not work
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).ToNot(BeNil())
		})

		// 15. testing if someone who is not the owener (invited) can not revoke other
		Specify("Testing to make sure someone who has been invited can not revoke access.", func() {
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			// Bob who is not the owner of the file is trying to revoke but should fail
			err = bob.RevokeAccess(bobFile, "alice")
			Expect(err).ToNot(BeNil())
		})

		// 16. testing to make sure someone with no access (someone who has been revoked) can not create new invitation
		Specify("Testing to make sure people with no access can not create new invitations.", func() {
			//ALice invites Bob
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			// Revokes bobs acccess poor bob
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			// Checking if bob cant load
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			// Checking to see if bob cant create new invitations (poor bob 2x)
			userlib.DebugMsg("Bob who has been revoked trying to create invitation (should not be able too),")
			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			_, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).ToNot(BeNil())
		})

		// 17. Checks if multi devices work (like overwriting reflects on different devices)
		// did not affect
		Specify("Seeing if Multi Device works (Overwiting on one device should show on the other Device).", func() {
			// Storing file on one device
			// aliceDesktop, err = client.InitUser("alice", defaultPassword)
			// Expect(err).To(BeNil())

			aliceDesktop := alice
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			// Logging on Different Device
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Laptop should be overwriting file.")
			err = aliceLaptop.StoreFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Now the overwriting should also show on the desktop login.")
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentThree)))
		})
	})

	Describe("security", func() {
		BeforeEach(func() {
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})

		// 18. Checking loadfile doesnt work for tamepred files
		// Nothing
		Specify("LoadFile should not work if the content has been tampered with.", func() {
			beforeKeys := make(map[uuid.UUID]bool)
			for k := range userlib.DatastoreGetMap() {
				beforeKeys[k] = true
			}

			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			// Find a new key and tamper with it
			dsAfter := userlib.DatastoreGetMap()
			found := false
			for k, v := range dsAfter {
				if !beforeKeys[k] {
					tamper := make([]byte, len(v))
					copy(tamper, v)
					tamper[0] ^= 0xFF
					userlib.DatastoreSet(k, tamper)
					found = true
					break
				}
			}
			Expect(found).To(BeTrue()) // sanity check: StoreFile wrote something

			_, err = alice.LoadFile(aliceFile)
			Expect(err).ToNot(BeNil())

		})

		// Specify("LoadFile should not work if the FileMeta has been tampered with.", func() {
		//     userlib.DebugMsg("storing file.")
		//     err = alice.StoreFile(aliceFile, []byte(contentOne))
		//     Expect(err).To(BeNil())

		//     // 1. Deterministically find the FileMeta UUID
		//     file_key, _ := userlib.HashKDF(alice.UserRootkey[:16], []byte("file-uuid"+aliceFile))
		//     fileMetaUUID, _ := uuid.FromBytes(file_key[:16])

		//     // 2. Get & Tamper
		//     metaBlob, ok := userlib.DatastoreGet(fileMetaUUID)
		//     Expect(ok).To(BeTrue())

		//     tamper := make([]byte, len(metaBlob))
		//     copy(tamper, metaBlob)
		//     tamper[0] ^= 0xFF
		//     userlib.DatastoreSet(fileMetaUUID, tamper)

		//     // 3. Check
		//     userlib.DebugMsg("loading tampered file.")
		//     _, err = alice.LoadFile(aliceFile)
		//     Expect(err).ToNot(BeNil())
		// })

		// 19. Checking acceptinivtation doesnt work if the invitation been tamepred with
		// Test 27
		Specify("accept invitation should not work if invitation has been tampered.", func() {
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())
			inviteUUID, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			inviteData, ok := userlib.DatastoreGet(inviteUUID)
			Expect(ok).To(BeTrue())
			tamper := make([]byte, len(inviteData))
			copy(tamper, inviteData)
			tamper[len(tamper)-1] ^= 0xFF //flips the last byte mac or sig
			userlib.DatastoreSet(inviteUUID, tamper)

			userlib.DebugMsg("bob accepting tamepred invitation.")
			err = bob.AcceptInvitation("alice", inviteUUID, bobFile)
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("random testing to see if we get coverage.", func() {
		var charles *client.User
		var charlesFile = "charlesFile.txt"

		BeforeEach(func() {
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())
		})

		Specify("user should be able to load file shared transitively like a to b to c", func() {
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			data, err := charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))
		})

		Specify("C loses access if B who invited C is revoked by A", func() {
			// return // waiting on revoke to be done
		})

		// Didnt work (maybe the previous test works already)
		Specify("Appending on one device sohuld show on the other device", func() {
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			aliceLaptop, err := client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			data, err := aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo)))
		})

		// Specify("loadfile should not work if sharedfile has been tampered", func() {
		// 	err = alice.StoreFile(aliceFile, []byte(contentOne))
		// 	Expect(err).To(BeNil())
		// 	invite, err := alice.CreateInvitation(aliceFile, "bob")
		// 	Expect(err).To(BeNil())
		// 	err = bob.AcceptInvitation("alice", invite, bobFile)
		// 	Expect(err).To(BeNil())

		// 	handleKey, err := userlib.HashKDF(bob.UserRootkey[:16], []byte("handle-uuid"+bobFile))
		// 	Expect(err).To(BeNil())
		// 	handleUUID, err := uuid.FromBytes(handleKey[:16])
		// 	Expect(err).To(BeNil())
		// 	handleBlob, ok := userlib.DatastoreGet(handleUUID)
		// 	Expect(ok).To(BeTrue())
		// 	tamper := make([]byte, len(handleBlob))
		// 	copy(tamper, handleBlob)
		// 	tamper[0] ^= 0xFF
		// 	userlib.DatastoreSet(handleUUID, tamper)
		// 	_, err = bob.LoadFile(bobFile)
		// 	Expect(err).ToNot(BeNil())

	})

	Describe("random testing", func() {
		BeforeEach(func() {
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())
		})

		// No testing
		Specify("multiple appends", func() {
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())

			expected := contentOne + contentTwo + contentThree
			Expect(string(data)).To(Equal(expected))
		})

		// doesnt pass any test
		Specify("testing", func() {
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			alice = nil

			alice, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(string(data)).To(Equal(contentOne))

		})

		// test 23
		Specify("Checking if username are case sensitive", func() {
			userlib.DebugMsg("Initializing user 'Alice' (capital A)")
			aliceUpper, err := client.InitUser("Alice", defaultPassword)
			Expect(err).To(BeNil())
			Expect(aliceUpper).ToNot(BeNil())

			userlib.DebugMsg("Getting user 'alice' (lowercase a)")
			aliceLower, err := client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			Expect(aliceUpper.Username).To(Equal("Alice"))
			Expect(aliceLower.Username).To(Equal("alice"))
		})

		// Test 22
		Specify("Handling larger files", func() {
			largeContent := make([]byte, 10000)
			for i := range largeContent { // filling with dummy (trash)
				largeContent[i] = 'A'
			}

			err = alice.StoreFile("large.bin", largeContent)
			Expect(err).To(BeNil())

			data, err := alice.LoadFile("large.bin")
			Expect(err).To(BeNil())
			Expect(data).To(Equal(largeContent))
		})

		// TEst 8
		Specify("Alice and Bob should be able to append to the same file efficiently", func() {
			err = alice.StoreFile(aliceFile, []byte("aliceone"))
			Expect(err).To(BeNil())

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			err = alice.AppendToFile(aliceFile, []byte("alicetwo"))
			Expect(err).To(BeNil())

			err = bob.AppendToFile(bobFile, []byte("bobone"))
			Expect(err).To(BeNil())

			err = alice.AppendToFile(aliceFile, []byte("alicethree"))
			Expect(err).To(BeNil())

			expected := "aliceonealicetwobobonealicethree"
			data, err := charles.LoadFile(aliceFile)
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(string(data)).To(Equal(expected))
		})

	})

	Describe("Start from 0 I hate this Project.", func() {
		BeforeEach(func() {
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())
		})

		// No Test
		Specify("Storing, loading, and appending to a empty file name", func() {
			filename := emptyString // basically "" but given as emptystring

			userlib.DebugMsg("storing file with empty name")
			err = alice.StoreFile(filename, []byte(contentOne))
			Expect(err).To(BeNil())

			data, err := alice.LoadFile(filename)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			err = alice.AppendToFile(filename, []byte(contentTwo))
			Expect(err).To(BeNil())

			data, err = alice.LoadFile(filename)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo)))
		})
		// nothing
		Specify("Overwiting ffile with identical content", func() {
			initial := []byte(contentOne)

			err = alice.StoreFile(aliceFile, initial)
			Expect(err).To(BeNil())

			err = alice.StoreFile(aliceFile, initial)
			Expect(err).To(BeNil())

			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal(initial))
		})

		// Nothing
		Specify("making sure wrong recipient can not accept invitation", func() {
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			inviteUUID, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("alice", inviteUUID, charlesFile)
			Expect(err).ToNot(BeNil())
		})

		// Test 12 and 25
		Specify("fake accepting when not been invited", func() {
			err = charles.StoreFile(charlesFile, []byte(contentOne))
			Expect(err).To(BeNil())

			realInvitation, err := charles.CreateInvitation(charlesFile, "bob")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("alice", realInvitation, "fake.txt")
			Expect(err).ToNot(BeNil())
		})

		// nothing
		Specify("can not open file if you been invited but did not accept invitation", func() {
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			_, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

		})

		// nothing
		Specify("if filenamed change the person invited should stillb eableto load", func() {
			aliceFilename := "aliceOne"
			bobFilename := "bobOne"

			err = alice.StoreFile(aliceFilename, []byte(contentOne))
			Expect(err).To(BeNil())

			invite, err := alice.CreateInvitation(aliceFilename, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFilename)
			Expect(err).To(BeNil())

			data, err := bob.LoadFile(bobFilename)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			_, err = bob.LoadFile(aliceFilename)
			Expect(err).ToNot(BeNil())
		})

		// nothing
		Specify("accepting invitation should not work if username is incorrect", func() {
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())
			err = bob.AcceptInvitation("charles", invite, bobFile)
			Expect(err).ToNot(BeNil())
		})

		// nothing
		Specify("no double accepting invitation", func() {
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).ToNot(BeNil())
		})

		//nothing
		// Specify("can not invite themselve", func() {
		// 	err = alice.StoreFile(aliceFile, []byte(contentOne))
		// 	Expect(err).To(BeNil())
		// 	_, err = alice.CreateInvitation(aliceFile, "alice")
		// 	Expect(err).ToNot(BeNil())
		// })

		// nothing
		Specify("multiple shares", func() {
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			inviteOne, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			inviteTwo, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", inviteOne, "bobOne")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", inviteTwo, "bobTwo")
			Expect(err).To(BeNil())
			dataOne, err := bob.LoadFile("bobOne")
			Expect(err).To(BeNil())
			Expect(dataOne).To(Equal([]byte(contentOne)))

			dataTwo, err := bob.LoadFile("bobTwo")
			Expect(err).To(BeNil())
			Expect(dataTwo).To(Equal([]byte(contentOne)))

			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			dataOne, _ = bob.LoadFile("bobOne")
			dataTwo, _ = bob.LoadFile("bobTwo")

			Expect(dataOne).To(Equal([]byte(contentOne + contentTwo)))
			Expect(dataTwo).To(Equal(dataOne))
		})

		// tnohing
		// Specify("sharing cycle should work without infinite recursion", func() {
		// 	err = alice.StoreFile("fileOne", []byte(contentOne))
		// 	Expect(err).To(BeNil())

		// 	inviteOne, err := alice.CreateInvitation("fileOne", "bob")
		// 	Expect(err).To(BeNil())

		// 	err = bob.AcceptInvitation("alice", inviteOne, "fileTwo")
		// 	Expect(err).To(BeNil())

		// 	inviteTwo, err := bob.CreateInvitation("fileTwo", "alice")
		// 	Expect(err).To(BeNil())

		// 	err = alice.AcceptInvitation("bob", inviteTwo, "fileThree")
		// 	Expect(err).To(BeNil())

		// 	data, err := alice.LoadFile("fileThree")
		// 	Expect(err).To(BeNil())
		// 	Expect(data).To(Equal([]byte(contentOne)))

		// 	err = alice.StoreFile("fileOne", []byte(contentTwo))
		// 	Expect(err).To(BeNil())

		// 	data, err = alice.LoadFile("fileThree")
		// 	Expect(err).To(BeNil())
		// 	Expect(data).To(Equal([]byte(contentTwo)))
		// })

		Specify("cant access file if revoke deven after logging out and back in", func() {
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			data, err := bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			bob = nil
			bobIsBack, err := client.GetUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			_, err = bobIsBack.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())
		})

		// nothing
		Specify("Alice should be able to reinvite after revoking", func() {
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			newInvite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			bobNewFile := "bob_is_back.txt"
			err = bob.AcceptInvitation("alice", newInvite, bobNewFile)
			Expect(err).To(BeNil())

			data, err := bob.LoadFile(bobNewFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))
		})

		// One of them is test 33
		Specify("revoking before accepting", func() {
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).ToNot(BeNil())
		})

		// One of them is test 33
		Specify("can not revoke own access unless youre the owner", func() {
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			err = bob.RevokeAccess(bobFile, "bob")
			Expect(err).ToNot(BeNil())
		})

		// nothing
		// Specify("whole revoke cycle thigny", func() {
		// 	err = alice.StoreFile(aliceFile, []byte(contentOne))
		// 	Expect(err).To(BeNil())

		// 	inviteOne, err := alice.CreateInvitation(aliceFile, "bob")
		// 	Expect(err).To(BeNil())

		// 	err = bob.AcceptInvitation("alice", inviteOne, bobFile)
		// 	Expect(err).To(BeNil())

		// 	inviteTwo, err := bob.CreateInvitation(bobFile, "charles")
		// 	Expect(err).To(BeNil())

		// 	err = alice.RevokeAccess(aliceFile, "bob")
		// 	Expect(err).To(BeNil())

		// 	err = charles.AcceptInvitation("bob", inviteTwo, charlesFile)
		// 	if err == nil {
		// 		_, err = charles.LoadFile(charlesFile)
		// 		Expect(err).ToNot(BeNil())
		// 	}
		// })

		// nothing
		// Specify("revoking shenangian", func() {
		// 	err = alice.StoreFile(aliceFile, []byte(contentOne))
		// 	Expect(err).To(BeNil())

		// 	inviteAliceToBob, err := alice.CreateInvitation(aliceFile, "bob")
		// 	Expect(err).To(BeNil())

		// 	err = bob.AcceptInvitation("alice", inviteAliceToBob, bobFile)
		// 	Expect(err).To(BeNil())

		// 	inviteAliceToCharles, err := alice.CreateInvitation(aliceFile, "charles")
		// 	Expect(err).To(BeNil())

		// 	err = charles.AcceptInvitation("alice", inviteAliceToCharles, "charles")
		// 	Expect(err).To(BeNil())

		// 	inviteBobToCharles, err := bob.CreateInvitation(bobFile, "charles")
		// 	Expect(err).To(BeNil())

		// 	err = charles.AcceptInvitation("bob", inviteBobToCharles, "charlestobob.txt")
		// 	Expect(err).To(BeNil())

		// 	err = alice.RevokeAccess(aliceFile, "bob")
		// 	Expect(err).To(BeNil())

		// 	data, err := charles.LoadFile("charlesshoulwork.txt")
		// 	Expect(err).To(BeNil())
		// 	Expect(data).To(Equal([]byte(contentOne)))

		// 	_, err = charles.LoadFile("charlestobob.txt")
		// 	Expect(err).ToNot(BeNil())
		// })

		// nothing
		Specify("Append should be able to be efficient even if file is beeeg", func() {
			err = alice.StoreFile(aliceFile, []byte("start"))
			Expect(err).To(BeNil())

			beegData := make([]byte, 1024)
			for i := 0; i < 50; i++ {
				err = alice.AppendToFile(aliceFile, beegData)
				Expect(err).To(BeNil())
			}

			err = alice.AppendToFile(aliceFile, []byte("end"))
			Expect(err).To(BeNil())

			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data[len(data)-3:]).To(Equal([]byte("end")))
		})

		// nothing
		Specify("Append should work with long file names.", func() {
			longName := ""
			for i := 0; i < 200; i++ {
				longName += "longname123"

			}

			err = alice.StoreFile(longName, []byte(contentOne))
			Expect(err).To(BeNil())

			err = alice.AppendToFile(longName, []byte(contentTwo))
			Expect(err).To(BeNil())

			data, err := alice.LoadFile(longName)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo)))
		})

		Specify("storing and loading should work with empty file name", func() {
			// emptyName := ""
			err = alice.StoreFile(emptyString, []byte(contentOne))
			Expect(err).To(BeNil())

			data, err := alice.LoadFile(emptyString)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))
		})

		// Specify("wait the opposite of the top one", func() {
		// 	emptyName := ""
		// 	err = alice.StoreFile(emptyName, []byte(contentOne))
		// 	Expect(err).ToNot(BeNil())

		// 	_, err = alice.LoadFile(emptyName)
		// 	Expect(err).ToNot(BeNil())

		// 	err = alice.AppendToFile(emptyName, []byte(contentTwo))
		// 	Expect(err).ToNot(BeNil())
		// })

		Specify("appending zero bytes to empty filename", func() {
			// emptyName := ""
			err = alice.StoreFile(emptyString, []byte(contentOne))
			Expect(err).To(BeNil())

			err = alice.AppendToFile(emptyString, []byte(""))
			Expect(err).To(BeNil())

			data, err := alice.LoadFile(emptyString)
			Expect(err).To(BeNil())

			Expect(data).To(Equal([]byte(contentOne)))
		})

		Specify("appending to normal content to empty filename should work", func() {
			// emptyName := ""
			err = alice.StoreFile(emptyString, []byte(contentOne))
			Expect(err).To(BeNil())

			err = alice.AppendToFile(emptyString, []byte(contentTwo))
			Expect(err).To(BeNil())

			data, err := alice.LoadFile(emptyString)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo)))
		})

		Specify("overwriting a file with same content ", func() {
			initial := []byte("some-data")
			err = alice.StoreFile(aliceFile, initial)
			Expect(err).To(BeNil())

			err = alice.StoreFile(aliceFile, initial)
			Expect(err).To(BeNil())

			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal(initial))
		})

		Specify("load and append should error if file no exist, func() {}", func() {
			no_exist := "ghost_scary.txt"
			_, err = alice.LoadFile(no_exist)
			Expect(err).ToNot(BeNil())

			err = alice.AppendToFile(no_exist, []byte("boo!"))
			Expect(err).ToNot(BeNil())
		})

		Specify("Users should be able to initialize and log in with an empty password", func() {
			u, err := client.InitUser("empty_pass_user", "")
			Expect(err).To(BeNil())
			Expect(u).ToNot(BeNil())

			u2, err := client.GetUser("empty_pass_user", "")
			Expect(err).To(BeNil())
			Expect(u2).ToNot(BeNil())

			_, err = client.GetUser("empty_pass_user", "password")
			Expect(err).ToNot(BeNil())
		})

		Specify("empty append to empty filename", func() {
			// emptyName := ""
			err = alice.StoreFile(emptyString, []byte(contentOne))
			Expect(err).To(BeNil())

			data, err := alice.LoadFile(emptyString)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))
			err = alice.AppendToFile(emptyString, []byte(""))
			Expect(err).To(BeNil())

			data, err = alice.LoadFile(emptyString)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))
		})

		Specify("appending to a file with empty name", func() {
			// emptyName := ""
			err = alice.StoreFile(emptyString, []byte(contentOne))
			Expect(err).To(BeNil())

			err = alice.AppendToFile(emptyString, []byte(contentTwo))
			Expect(err).To(BeNil())

			data, err := alice.LoadFile(emptyString)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo)))
		})

		Specify("appending nothing to empty file name error", func() {

			err = alice.AppendToFile(emptyString, []byte(""))
			Expect(err).ToNot(BeNil())

		})

		Specify("appending something to empty file name error:", func() {

			err = alice.AppendToFile(emptyString, []byte(contentTwo))

			Expect(err).ToNot(BeNil())
		})

		Specify("load and append to file that no exist", func() {
			noFile := "no_file.txt"
			_, err = alice.LoadFile(noFile)
			Expect(err).ToNot(BeNil())

			err = alice.AppendToFile(noFile, []byte(contentOne))
			Expect(err).ToNot(BeNil())
		})

		// test 17
		Specify("loooooooooooooooooooooooooooooooooooong appends", func() {
			err = alice.StoreFile(aliceFile, []byte("Start"))
			Expect(err).To(BeNil())
			expectedData := "Start"
			appendnew := "next"

			for i := 0; i < 100; i++ {
				err = alice.AppendToFile(aliceFile, []byte(appendnew))
				Expect(err).To(BeNil())
				expectedData += appendnew
			}

			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(string(data)).To(Equal(expectedData))
		})

		Specify("load append store all should work with empty name ", func() {
			err = alice.StoreFile(emptyString, []byte(contentOne))
			Expect(err).To(BeNil())

			data, err := alice.LoadFile(emptyString)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))
			err = alice.AppendToFile(emptyString, []byte(contentTwo))
			Expect(err).To(BeNil())

			data, err = alice.LoadFile(emptyString)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo)))
		})

		// SUPER LINE
		Specify("Tampering with User struct should cause GetUser to fail", func() {

			aliceUUID, err := uuid.FromBytes(userlib.Hash([]byte("user-" + "alice"))[:16])
			Expect(err).To(BeNil())

			userBlob, ok := userlib.DatastoreGet(aliceUUID)
			Expect(ok).To(BeTrue())

			// Tamper with the blob
			tampered := make([]byte, len(userBlob))
			copy(tampered, userBlob)
			if len(tampered) > 10 {
				tampered[10] ^= 0xFF
			} else {
				tampered[0] ^= 0xFF
			}
			userlib.DatastoreSet(aliceUUID, tampered)

			// Now GetUser should fail due to MAC mismatch / malformed blob
			_, err = client.GetUser("alice", defaultPassword)
			Expect(err).ToNot(BeNil())
		})

		Specify("Tampering with filenode should cause Load to go boooom", func() {
			aliceUUID, err := uuid.FromBytes(userlib.Hash([]byte("user-" + "alice"))[:16])
			Expect(err).To(BeNil())

			userBlob, ok := userlib.DatastoreGet(aliceUUID)
			Expect(ok).To(BeTrue())

			// Tamper with the blob
			tampered := make([]byte, len(userBlob))
			copy(tampered, userBlob)
			if len(tampered) > 10 {
				tampered[10] ^= 0xFF
			} else {
				tampered[0] ^= 0xFF
			}
			userlib.DatastoreSet(aliceUUID, tampered)

			// Now GetUser should fail due to MAC mismatch / malformed blob
			_, err = client.GetUser("alice", defaultPassword)
			Expect(err).ToNot(BeNil())
		})

		Specify("tampering any of the nodes in the linktlist should fail loadfile", func() {
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			ds := userlib.DatastoreGetMap() // getting all datastore entry

			count := 0
			for k, v := range ds {
				if count == 2 { // skipping first few to avoi the user struct
					tampered := make([]byte, len(v))
					copy(tampered, v)
					if len(tampered) > 5 {
						tampered[5] ^= 0xFF
					}

					userlib.DatastoreSet(k, tampered)
					break
				}
				count++
			}

			_, err = alice.LoadFile(aliceFile)
		})

		Specify("Tampering with the invitation should cause accept to fial", func() {
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			inviteUUID, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			inviteBlob, ok := userlib.DatastoreGet(inviteUUID)
			Expect(ok).To(BeTrue())

			tampered := make([]byte, len(inviteBlob))
			copy(tampered, inviteBlob)
			tampered[len(tampered)-1] ^= 0xFF
			userlib.DatastoreSet(inviteUUID, tampered)

			err = bob.AcceptInvitation("alice", inviteUUID, bobFile)
			Expect(err).ToNot(BeNil())
		})

		Specify("delting filenode causes load to fail", func() {

			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			// Snapshot keys before append
			beforeKeys := make(map[uuid.UUID]bool)
			for k := range userlib.DatastoreGetMap() {
				beforeKeys[k] = true
			}

			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			// Find a new key created by AppendToFile
			dsAfter := userlib.DatastoreGetMap()
			var newNodeUUID uuid.UUID
			found := false
			for k := range dsAfter {
				if !beforeKeys[k] {
					newNodeUUID = k
					found = true
					break
				}
			}
			Expect(found).To(BeTrue()) // sanity check: append wrote something

			userlib.DatastoreDelete(newNodeUUID)

			_, err = alice.LoadFile(aliceFile)
			Expect(err).ToNot(BeNil())
		})

		// Specify("tampering with the public key in keystore sohuld make createinv fail", func() {
		// 	err = alice.StoreFile(aliceFile, []byte(contentOne))
		// 	Expect(err).To(BeNil())

		// 	// Save original Bob public key
		// 	bobPubKey, ok := userlib.KeystoreGet("bob_pke")
		// 	Expect(ok).To(BeTrue())

		// 	// Replace with fake key
		// 	fakePubKey, _, err := userlib.PKEKeyGen()
		// 	Expect(err).To(BeNil())
		// 	err = userlib.KeystoreSet("bob_pke", fakePubKey)
		// 	Expect(err).To(BeNil())

		// 	// Create invite should still succeed
		// 	inviteUUID, err := alice.CreateInvitation(aliceFile, "bob")
		// 	Expect(err).To(BeNil())

		// 	// But Bob cannot accept it (decryption should fail somewhere)
		// 	err = bob.AcceptInvitation("alice", inviteUUID, bobFile)
		// 	Expect(err).ToNot(BeNil())

		// 	// Restore original key so other tests aren't affected
		// 	err = userlib.KeystoreSet("bob_pke", bobPubKey)
		// 	Expect(err).To(BeNil())
		// })

		// Specify("Tampering with signature verification key causes acceptinv to fail", func() {
		// 	err = alice.StoreFile(aliceFile, []byte(contentOne))
		// 	Expect(err).To(BeNil())

		// 	inviteUUID, err := alice.CreateInvitation(aliceFile, "bob")
		// 	Expect(err).To(BeNil())

		// 	aliceVerifyKey, ok := userlib.KeystoreGet("alice_dsv")
		// 	Expect(ok).To(BeTrue())

		// 	_, fakeVerifyKey, err := userlib.DSKeyGen()
		// 	Expect(err).To(BeNil())

		// 	err = userlib.KeystoreSet("alice_dsv", fakeVerifyKey)
		// 	Expect(err).To(BeNil())

		// 	err = bob.AcceptInvitation("alice", inviteUUID, bobFile)
		// 	Expect(err).ToNot(BeNil())

		// 	// AcceptInvitation should fail (signature won't verify)
		// 	err = bob.AcceptInvitation("alice", inviteUUID, bobFile)
		// 	Expect(err).ToNot(BeNil())

		// 	// Restore correct key
		// 	err = userlib.KeystoreSet("alice_dsv", aliceVerifyKey)
		// 	Expect(err).To(BeNil())
		// })

		Specify("Tampering with file data sohuld be detected by owner and the shared user", func() {

			beforeKeys := make(map[uuid.UUID]bool)
			for k := range userlib.DatastoreGetMap() {
				beforeKeys[k] = true
			}

			// 2. Alice stores the file
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			// 3. Find the new keys that StoreFile created (metadata + base node)
			afterStore := userlib.DatastoreGetMap()
			var fileKeys []uuid.UUID
			for k := range afterStore {
				if !beforeKeys[k] {
					fileKeys = append(fileKeys, k)
				}
			}
			Expect(len(fileKeys)).To(BeNumerically(">", 0)) // sanity: StoreFile wrote something

			// 4. Share to Bob and make sure both can load before tampering
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			_, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())

			_, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())

			// 5. Tamper all blobs that belonged to this file when it was created
			for _, k := range fileKeys {
				blob, ok := userlib.DatastoreGet(k)
				Expect(ok).To(BeTrue())

				if len(blob) == 0 {
					continue
				}
				tampered := make([]byte, len(blob))
				copy(tampered, blob)
				tampered[0] ^= 0xFF
				userlib.DatastoreSet(k, tampered)
			}

			// 6. Now both Alice and Bob should see an error when loading
			_, err = alice.LoadFile(aliceFile)
			Expect(err).ToNot(BeNil())

			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

		})

	})

	//------------------------------------Arman's tests--------------------------------------------------------

	// 1: checks MAC to see tampering
	// 2: checks for correct password
	// 3: checks for incorrect password ( same thing copied twice lmao)
	Describe("Arman's tests", func() {
		// Test 29
		// 1: checks MAC to see tampering
		Specify("Owner should detect tampered file metadata", func() {
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			metaKey, err := userlib.HashKDF(alice.UserRootkey[:16], []byte("file-uuid"+aliceFile))
			Expect(err).To(BeNil())
			metaUUID, err := uuid.FromBytes(metaKey[:16])
			Expect(err).To(BeNil())

			metaBlob, ok := userlib.DatastoreGet(metaUUID)
			Expect(ok).To(BeTrue())

			tampered := make([]byte, len(metaBlob))
			copy(tampered, metaBlob)
			tampered[0] ^= 0xFF
			userlib.DatastoreSet(metaUUID, tampered)

			_, err = alice.LoadFile(aliceFile)
			Expect(err).ToNot(BeNil())
		})

		Specify("Owner and shared user should detect tampered file node", func() {
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			// snapshot keys to find file blobs created by StoreFile
			before := make(map[uuid.UUID]bool)
			for k := range userlib.DatastoreGetMap() {
				before[k] = true
			}

			// append to force at least one node
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			after := userlib.DatastoreGetMap()
			var nodeUUID uuid.UUID
			for k := range after {
				if !before[k] {
					nodeUUID = k
					break
				}
			}
			Expect(nodeUUID).ToNot(Equal(uuid.Nil))

			nodeBlob, ok := userlib.DatastoreGet(nodeUUID)
			Expect(ok).To(BeTrue())

			tampered := make([]byte, len(nodeBlob))
			copy(tampered, nodeBlob)
			tampered[0] ^= 0xFF
			userlib.DatastoreSet(nodeUUID, tampered)

			_, err = alice.LoadFile(aliceFile)
			Expect(err).ToNot(BeNil())

			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())
		})

		// 2: checks for correct password
		Specify("GetUser should fail with the wrong password", func() {
			_, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			// wrong password should NOT decrypt/MAC-verify
			_, err = client.GetUser("alice", "not-the-password")
			Expect(err).ToNot(BeNil())
		})

		// 3: checks for incorrect password
		Specify("should fail if the wrong password is inputted", func() {

			_, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			_, err = client.GetUser("alice", "not_a_good_passkey")
			Expect(err).ToNot(BeNil())
		})

	})
})
