package client

// CS 161 Project 2

// Only the following imports are allowed! ANY additional imports
// may break the autograder!
// - bytes
// - encoding/hex
// - encoding/json
// - errors
// - fmt
// - github.com/cs161-staff/project2-userlib
// - github.com/google/uuid
// - strconv
// - strings

import (
	"encoding/json"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	// hex.EncodeToString(...) is useful for converting []byte to string

	// Useful for string manipulation

	// Useful for formatting strings (e.g. `fmt.Sprintf`).
	"fmt"

	// Useful for creating new error messages to return using errors.New("...")
	"errors"

	// Optional.
	_ "strconv"
)

// This serves two purposes: it shows you a few useful primitives,
// and suppresses warnings for imports not being used. It can be
// safely deleted!
func someUsefulThings() {

	// Creates a random UUID.
	randomUUID := uuid.New()

	// Prints the UUID as a string. %v prints the value in a default format.
	// See https://pkg.go.dev/fmt#hdr-Printing for all Golang format string flags.
	userlib.DebugMsg("Random UUID: %v", randomUUID.String())

	// Creates a UUID deterministically, from a sequence of bytes.
	hash := userlib.Hash([]byte("user-structs/alice"))
	deterministicUUID, err := uuid.FromBytes(hash[:16])
	if err != nil {
		// Normally, we would `return err` here. But, since this function doesn't return anything,
		// we can just panic to terminate execution. ALWAYS, ALWAYS, ALWAYS check for errors! Your
		// code should have hundreds of "if err != nil { return err }" statements by the end of this
		// project. You probably want to avoid using panic statements in your own code.
		panic(errors.New("An error occurred while generating a UUID: " + err.Error()))
	}
	userlib.DebugMsg("Deterministic UUID: %v", deterministicUUID.String())

	// Declares a Course struct type, creates an instance of it, and marshals it into JSON.
	type Course struct {
		name      string
		professor []byte
	}

	course := Course{"CS 161", []byte("Nicholas Weaver")}
	courseBytes, err := json.Marshal(course)
	if err != nil {
		panic(err)
	}

	userlib.DebugMsg("Struct: %v", course)
	userlib.DebugMsg("JSON Data: %v", courseBytes)

	// Generate a random private/public keypair.
	// The "_" indicates that we don't check for the error case here.
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("PKE Key Pair: (%v, %v)", pk, sk)

	// Here's an example of how to use HBKDF to generate a new key from an input key.
	// Tip: generate a new key everywhere you possibly can! It's easier to generate new keys on the fly
	// instead of trying to think about all of the ways a key reuse attack could be performed. It's also easier to
	// store one key and derive multiple keys from that one key, rather than
	originalKey := userlib.RandomBytes(16)
	derivedKey, err := userlib.HashKDF(originalKey, []byte("mac-key"))
	if err != nil {
		panic(err)
	}
	userlib.DebugMsg("Original Key: %v", originalKey)
	userlib.DebugMsg("Derived Key: %v", derivedKey)

	// A couple of tips on converting between string and []byte:
	// To convert from string to []byte, use []byte("some-string-here")
	// To convert from []byte to string for debugging, use fmt.Sprintf("hello world: %s", some_byte_arr).
	// To convert from []byte to string for use in a hashmap, use hex.EncodeToString(some_byte_arr).
	// When frequently converting between []byte and string, just marshal and unmarshal the data.
	//
	// Read more: https://go.dev/blog/strings

	// Here's an example of string interpolation!
	_ = fmt.Sprintf("%s_%d", "file", 1)
}

// This is the type definition for the User struct.
// A Go struct is like a Python or Java class - it can have attributes
// (e.g. like the Username attribute) and methods (e.g. like the StoreFile method below).

// for this project my user struct is gonna be MAD VIBIN ,Bossin bossin 360 no scope that shit
// like the good old days not thhis bullshit of an excuse they call bo7
// what happened to making good games? why did video games die?
// I remember the golden age of cod with MW and all the other shit but now this is just pure otter BS
// and dont get me started on the matchmaking. cause then I will kill myself on the spot.
// anyways enough and lets get back into this
// the user struct holds the username, the the persons kie khar, the persons private keys and their private keys, both enc and MAC ones
// N is a random number number I was planning on using for creating random keys and so each user has their own N
// but I later on realized I dont need it but I still put it there
// I used a User ROOT key instead, same thing but blah blah
type User struct {
	Username         string
	PrivEncKey       userlib.PKEDecKey
	PubEncKey        userlib.PKEEncKey
	PrivSignatureKey userlib.DSSignKey
	PubSignatureKey  userlib.DSVerifyKey
	N                []byte
	FilesMapUUID     uuid.UUID
	UserRootkey      []byte

	// You can add other attributes here if you want! But note that in order for attributes to
	// be included when this struct is serialized to/from JSON, they must be capitalized.
	// On the flipside, if you have an attribute that you want to be able to access from
	// this struct's methods, but you DON'T want that value to be included in the serialized value
	// of this struct that's stored in datastore, then you can use a "private" variable (e.g. one that
	// begins with a lowercase letter).
}
type FileMeta struct {
	Owner        string
	FileName     string
	EncKey       []byte
	MACKey       []byte
	Base_Node    uuid.UUID
	Current_Node uuid.UUID
	Owner_map    uuid.UUID
}

// im using a linked list thingy for fast appends so its a linked list thingy happily ever after
type FileNode struct {
	Content []byte
	Next    uuid.UUID //or a UUID not sure yet about that

}

type Invitation struct {
	OGOwner      string
	FileMetaData uuid.UUID // the address of the metadatafile struct
	MacKey       []byte
	EncKey       []byte
	Revoked      bool
}

type SharedFile struct {
	OGOwner     string
	NodeAddress uuid.UUID
	InviteUUID  uuid.UUID // the address of the invitation struct in which we would retreieve the live passwords from
	// Owner_map   uuid.UUID //  dont need this anymore thank gawd
	//
	// the uuid of the hashmap which would store the the people in which the file has been shared with
	// could be an array or a hashmap, havnt decided yet on which one

}

// so in my testings I ran into an issue of my invite strcuts being waaaaaaaaaayyyy tooo long
// and so I created this pack invitation helpers to help me "PACK" my stuff
// also I didnt come up with this section myself,
// I used pookie GPT to for this part since course staff was being real funny in helping for this project

func packInvitation(inv Invitation) ([]byte, error) {
	if len(inv.EncKey) < 16 || len(inv.MacKey) < 16 {
		return nil, errors.New("invitation keys are aaaaaaaaaaaaaaaaaa waaaay too short typr shiiii")
	}

	uuidBytes, err := inv.FileMetaData.MarshalBinary()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 0, 49)

	buf = append(buf, inv.EncKey[:16]...)

	buf = append(buf, inv.MacKey[:16]...)

	buf = append(buf, uuidBytes...)

	if inv.Revoked {
		buf = append(buf, byte(1))
	} else {
		buf = append(buf, byte(0))
	}

	return buf, nil
}

func unpackInvitation(data []byte) (Invitation, error) {
	if len(data) != 49 {
		return Invitation{}, errors.New("bad invitation size")
	}

	inv := Invitation{}
	inv.EncKey = make([]byte, 16)
	inv.MacKey = make([]byte, 16)

	copy(inv.EncKey, data[0:16])
	copy(inv.MacKey, data[16:32])

	var id uuid.UUID
	if err := id.UnmarshalBinary(data[32:48]); err != nil {
		return Invitation{}, err
	}
	inv.FileMetaData = id

	inv.Revoked = (data[48] == 1)

	// we dont have to encode OGowner is im not really using it ( skull emoji)

	return inv, nil
}

// some wise words before we start coding, ( also remeber to fix that thing with the invitation pointer not
// not having a live last node) lets derive a secret location UUID from the password or from the passwoerd driven key
// and then store the enctpyed user and MACS there

// NOTE: The following methods have toy (insecure!) implementations.

//helper function for generating random shit for N:

//check if user exists or not by deriving their datastore UUID from the password if something exists then we return an error
// after that we generate 2 keypairs and store the public ones on the the keystore and store the private ones locally
// make a random salt (or N) for the argon2key function so we generate our source key from the password.
//split that base key into a userenckey and a userMackey and potentially a userlocation key
//build a user struct, then encrpyt and mac the user struct, store it under a determenistic UUID in the datastore thingy
//

func InitUser(username string, password string) (userdataptr *User, err error) {
	//if someone dosnt input a username then we error:
	if username == "" {
		return nil, errors.New("no username was provided")
	}
	salt := userlib.RandomBytes(16)
	base_Key := userlib.Argon2Key([]byte(password), salt, 32)
	//here I will split the base key into the 3 keys I mentioned
	EncKey, err := userlib.HashKDF(base_Key[:16], []byte("user-enc"))
	if err != nil {
		return nil, err
	}
	Mac_Key, err := userlib.HashKDF(base_Key[:16], []byte("user-mac"))
	if err != nil {
		return nil, err
	}
	Root_key, err := userlib.HashKDF(base_Key[:16], []byte("user-root"))
	if err != nil {
		return nil, err
	}
	userUUID, err := uuid.FromBytes(userlib.Hash([]byte("user-" + username))[:16])
	if err != nil {
		return nil, err
	}

	//dont forget to check if the user already exists or not, if the user already exists we must return an error
	if _, ok := userlib.DatastoreGet(userUUID); ok {
		return nil, errors.New("user already exists so pick another username or something")
	}
	// now we genarate the key pairs we talked about earlier and in the spec
	PubEnc, PrivEnc, err := userlib.PKEKeyGen()
	if err != nil {
		return nil, err
	}

	privSig, pubSig, err := userlib.DSKeyGen()
	if err != nil {
		return nil, err
	}
	// now we put our keys in the keystore thingy
	err = userlib.KeystoreSet(username+"_pke", PubEnc)
	if err != nil {
		return nil, err

	}
	err = userlib.KeystoreSet(username+"_dsv", pubSig)
	if err != nil {
		return nil, err
	}

	//now we build the user object

	userobj := User{
		Username:         username,
		PrivEncKey:       PrivEnc,
		PubEncKey:        PubEnc,
		PrivSignatureKey: privSig,
		PubSignatureKey:  pubSig,
		N:                userlib.RandomBytes(16), // or convert from RandomBytes if RandomInt not present
		FilesMapUUID:     uuid.New(),
		UserRootkey:      Root_key,
	}
	filesMap := make(map[string]uuid.UUID)
	fmBytes, err := json.Marshal(filesMap)

	if err != nil {
		return nil, err
	}
	enc_FM := userlib.SymEnc(EncKey[:16], userlib.RandomBytes(16), fmBytes)
	fm_Tag, err := userlib.HMACEval(Mac_Key[:16], enc_FM)
	if err != nil {
		return nil, err
	}
	files_blob := struct {
		C   []byte
		Tag []byte
	}{C: enc_FM, Tag: fm_Tag}
	b, _ := json.Marshal(files_blob)
	userlib.DatastoreSet(userobj.FilesMapUUID, b)

	userBytes, err := json.Marshal(userobj)
	if err != nil {
		return nil, err
	}

	// now that we seriliazed the user struct, we will try to encrpyt it and then store it on data store
	iv := userlib.RandomBytes(16)
	encrypted_text := userlib.SymEnc(EncKey[:16], iv, userBytes)
	//now we will create the MAC
	mac_input, err := json.Marshal(struct {
		Salt []byte
		IV   []byte
		C    []byte
	}{
		Salt: salt,
		IV:   iv,
		C:    encrypted_text,
	})
	if err != nil {
		return nil, err
	}
	tag, err := userlib.HMACEval(Mac_Key[:16], mac_input)
	if err != nil {
		return nil, err
	}
	storeBlob, err := json.Marshal(struct {
		Salt []byte
		IV   []byte
		C    []byte
		Tag  []byte
	}{
		Salt: salt,
		IV:   iv,
		C:    encrypted_text,
		Tag:  tag,
	})
	if err != nil {
		return nil, err
	}

	userlib.DatastoreSet(userUUID, storeBlob)

	return &userobj, nil
	//  return &userdata, nil
}

// look up the usernameUUID and see if it exists, extract the salt, derive the keys using argon2key and hashkdf
//
/// then decrypt the userblob and unmarhsal it

// return the userdataptr or an error
func GetUser(username string, password string) (userdataptr *User, err error) {
	userUUID, err := uuid.FromBytes(userlib.Hash([]byte("user-" + username))[:16])
	if err != nil {
		return nil, err
	}
	userBlob, ok := userlib.DatastoreGet(userUUID)

	if !ok {
		return nil, errors.New("bro think you may be fucked, there are no known users with such a name")
	}

	var locator struct {
		Salt []byte
		IV   []byte
		C    []byte
		Tag  []byte
	}
	if err := json.Unmarshal(userBlob, &locator); err != nil {
		return nil, errors.New("hmmm something is wrong, not sure what tho so you maybe fucked")
	}

	base_kilid := userlib.Argon2Key([]byte(password), locator.Salt, 32)

	encKey, err := userlib.HashKDF(base_kilid[:16], []byte("user-enc"))
	if err != nil {
		return nil, err
	}
	macKey, err := userlib.HashKDF(base_kilid[:16], []byte("user-mac"))
	if err != nil {
		return nil, err
	}

	if len(locator.C) == 0 || len(locator.IV) == 0 || len(locator.Tag) == 0 {

		return nil, errors.New("encrypted user blob missing from locator record")
	}

	macInput, err := json.Marshal(struct {
		Salt []byte
		IV   []byte
		C    []byte
	}{
		Salt: locator.Salt,
		IV:   locator.IV,
		C:    locator.C,
	})
	if err != nil {
		return nil, err
	}

	expected_tag, err := userlib.HMACEval(macKey[:16], macInput)
	if err != nil {
		return nil, err
	}
	if !userlib.HMACEqual(expected_tag, locator.Tag) {
		return nil, errors.New("user blob MAC mismatch (tampering or wrong password)")
	}
	plain := userlib.SymDec(encKey[:16], locator.C)

	var userobj User
	if err := json.Unmarshal(plain, &userobj); err != nil {
		return nil, err
	}

	return &userobj, nil
}

// 	var userdata User
// 	userdataptr = &userdata
// 	return userdataptr, nil
// }

func (userdata *User) StoreFile(filename string, content []byte) (err error) {

	// get the key for the deterministic UUID, we are using the root key and the hashkdf function
	// to derive a key for the UUID ( just the "ID" for now)
	file_key, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("file-uuid"+filename))

	if err != nil {
		return err
	}
	metadata_enc_key, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("file-meta-enc"+filename))
	if err != nil {
		return err
	}
	metadata_mac_key, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("file-meta-mac"+filename))
	if err != nil {
		return err
	}

	//here we generate the UUID form the key we derived earlier
	FileMetaUUID, err := uuid.FromBytes(file_key[:16])
	if err != nil {
		return err
	}

	makeNode := func(encKey []byte, macKey []byte) (nodeUUID uuid.UUID, err error) {
		nodeUUID = uuid.New()
		node := FileNode{
			Content: content,
			Next:    uuid.Nil,
		}
		nodeBytes, err := json.Marshal(node)
		if err != nil {
			return uuid.Nil, err
		}
		nodeCT := userlib.SymEnc(encKey[:16], userlib.RandomBytes(16), nodeBytes)
		nodeTag, err := userlib.HMACEval(macKey[:16], nodeCT)
		if err != nil {
			return uuid.Nil, err
		}
		nodeBlob, err := json.Marshal(struct {
			C   []byte
			Tag []byte
		}{
			C:   nodeCT,
			Tag: nodeTag,
		})
		if err != nil {
			return uuid.Nil, err
		}
		userlib.DatastoreSet(nodeUUID, nodeBlob)
		return nodeUUID, nil
	}

	// we get the file blob, if nothing exists then we create it and store, if something
	// already exists then we will replace it with the new stuff
	metablob, ok := userlib.DatastoreGet(FileMetaUUID)
	if !ok {
		// generate a pair of mac and enc keys to store in the filemetadata struct and to encrpyt the nodes
		fileEncKey := userlib.RandomBytes(16)
		fileMacKey := userlib.RandomBytes(16)

		//then we try to create the first node
		nodeuuid, err := makeNode(fileEncKey, fileMacKey)
		if err != nil {

			return err
		}

		// and then the filemetastruct e difined earlier it comes to life here for each store file CALL
		file_meta := FileMeta{
			Owner:        userdata.Username,
			FileName:     filename,
			EncKey:       fileEncKey,
			MACKey:       fileMacKey,
			Base_Node:    nodeuuid,
			Current_Node: nodeuuid,
		}
		//then we serialize it
		metabytes, err := json.Marshal(file_meta)
		if err != nil {
			return err
		}
		// then we encrpyt it
		meta_cypher_text := userlib.SymEnc(metadata_enc_key[:16], userlib.RandomBytes(16), metabytes)
		meta_tag, err := userlib.HMACEval(metadata_mac_key[:16], meta_cypher_text)
		if err != nil {
			return err
		}

		finalMeta, err := json.Marshal(struct {
			C   []byte
			Tag []byte
		}{
			C:   meta_cypher_text,
			Tag: meta_tag,
		})
		if err != nil {
			return err
		}

		userlib.DatastoreSet(FileMetaUUID, finalMeta)
		return nil

	}
	// now if the thing alreadyexists we fetch it an
	//d renamee it
	var wrapped struct {
		C   []byte
		Tag []byte
	}
	//then we grab it and de un marshal the shit out of that thing
	if err := json.Unmarshal(metablob, &wrapped); err != nil {
		return err
	}

	//then we check the fuck out of the tag of that thingy

	tag_check, err := userlib.HMACEval(metadata_mac_key[:16], wrapped.C)
	if err != nil {
		return err
	}
	// we error if someone tried touching the files
	if !userlib.HMACEqual(tag_check, wrapped.Tag) {
		return errors.New("oh oh some is touching the files")
	}

	//now we decrpyt

	meta_plain := userlib.SymDec(metadata_enc_key[:16], wrapped.C)
	var meta FileMeta
	if err := json.Unmarshal(meta_plain, &meta); err != nil {
		return err
	}

	// we can reuse existing per-file keys (meta.EncKey, meta.MACKey)
	nodeUUID, err := makeNode(meta.EncKey, meta.MACKey)
	if err != nil {
		return err
	}

	// update metadata to point to new content
	// overwrite semantics: base = current node
	meta.Base_Node = nodeUUID
	meta.Current_Node = nodeUUID

	meta_bytes, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	meta_cypher_text := userlib.SymEnc(metadata_enc_key[:16], userlib.RandomBytes(16), meta_bytes)
	meta_tag, err := userlib.HMACEval(metadata_mac_key[:16], meta_cypher_text)
	if err != nil {
		return err
	}

	finalMeta, err := json.Marshal(struct {
		C   []byte
		Tag []byte
	}{
		C:   meta_cypher_text,
		Tag: meta_tag,
	})
	if err != nil {
		return err
	}

	userlib.DatastoreSet(FileMetaUUID, finalMeta)
	return nil
}

// 	fmblob, ok := userlib.DatastoreGet(userdata.FilesMapUUID)
// 	if !ok {
// 		return errors.New("uh filemap was never initialized chief")
// 	}
// 	var fm_wrapped struct {
// 		C   []byte
// 		Tag []byte
// 	}
// 	if err := json.Unmarshal(fmblob, &fm_wrapped); err != nil {
// 		return err
// 	}

// 	storageKey, err := uuid.FromBytes(userlib.Hash([]byte(filename + userdata.Username))[:16])
// 	if err != nil {
// 		return err
// 	}
// 	contentBytes, err := json.Marshal(content)
// 	if err != nil {
// 		return err
// 	}
// 	userlib.DatastoreSet(storageKey, contentBytes)
// 	return
// }

//--------------------------------------------MY APPEND TO FILE---------------------------------------------------------

func (userdata *User) AppendToFile(filename string, content []byte) error {
	// if filename == "" {
	// 	return errors.New("empty filename")
	// }

	//
	//  Try OWNER path first
	//
	// Use the SAME derivations as in StoreFile / CreateInvitation.
	fileKey, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("file-uuid"+filename))
	if err != nil {
		return err
	}
	fileMetaUUID, err := uuid.FromBytes(fileKey[:16])
	if err != nil {
		return err
	}

	metaEncMaterial, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("file-meta-enc"+filename))
	if err != nil {
		return err
	}
	metaMacMaterial, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("file-meta-mac"+filename))
	if err != nil {
		return err
	}
	metaEncKey := metaEncMaterial[:16]
	metaMacKey := metaMacMaterial[:16]

	metaBlob, ok := userlib.DatastoreGet(fileMetaUUID)
	if ok {
		// We found something at the owner's FileMeta location. Verify integrity.
		var wrapped struct {
			C   []byte
			Tag []byte
		}
		//then we unmarshal like before
		if err := json.Unmarshal(metaBlob, &wrapped); err != nil {
			return err
		}
		// 			then we check the tag like before
		tagCheck, err := userlib.HMACEval(metaMacKey, wrapped.C)
		if err != nil {
			return err
		}
		if !userlib.HMACEqual(tagCheck, wrapped.Tag) {
			return errors.New("file metadata tampered")
		}
		// then we decrpty like the other time
		metaP_lain := userlib.SymDec(metaEncKey, wrapped.C)

		var meta FileMeta
		if err := json.Unmarshal(metaP_lain, &meta); err != nil {
			return err
		}

		// Owner has the per-file keys and the current tail node.
		fileEncKey := meta.EncKey
		fileMacKey := meta.MACKey

		// Helper to create and store a new node with given content.
		makeNode := func(data []byte) (uuid.UUID, error) {
			nodeUUID := uuid.New()
			node := FileNode{
				Content: data,
				Next:    uuid.Nil,
			}
			nodeBytes, err := json.Marshal(node)
			if err != nil {
				return uuid.Nil, err
			}
			ct := userlib.SymEnc(fileEncKey[:16], userlib.RandomBytes(16), nodeBytes)
			tag, err := userlib.HMACEval(fileMacKey[:16], ct)
			if err != nil {
				return uuid.Nil, err
			}
			nodeBlob, err := json.Marshal(struct {
				C   []byte
				Tag []byte
			}{
				C:   ct,
				Tag: tag,
			})
			if err != nil {
				return uuid.Nil, err
			}
			userlib.DatastoreSet(nodeUUID, nodeBlob)
			return nodeUUID, nil
		}

		// create the new tail node for the appended content.
		newNodeUUID, err := makeNode(content)
		if err != nil {
			return err
		}

		//moew the real tail by walking from Base_Node.
		// Non-owners can append without updating FileMeta.Current_Node,
		// so we *cannot* trust meta.Current_Node to be the actual tail.
		if meta.Base_Node == uuid.Nil {
			// No nodes yet (somehow) – treat this as first node.
			meta.Base_Node = newNodeUUID
			meta.Current_Node = newNodeUUID
		} else {
			// Walk from Base_Node to the true tail.
			tailUUID := meta.Base_Node
			var tailNode FileNode

			for {
				tailBlob, ok := userlib.DatastoreGet(tailUUID)
				if !ok {
					return errors.New("tail node missing from datastore")
				}

				var tailWrapped struct {
					C   []byte
					Tag []byte
				}
				if err := json.Unmarshal(tailBlob, &tailWrapped); err != nil {
					return err
				}

				tailTagCheck, err := userlib.HMACEval(fileMacKey[:16], tailWrapped.C)
				if err != nil {
					return err
				}
				if !userlib.HMACEqual(tailTagCheck, tailWrapped.Tag) {
					return errors.New("tail node tampered")
				}

				tailP_lain := userlib.SymDec(fileEncKey[:16], tailWrapped.C)
				if err := json.Unmarshal(tailP_lain, &tailNode); err != nil {
					return err
				}
				//AHA we got the tail here
				if tailNode.Next == uuid.Nil {
					// Found actual tail.
					break
				}
				tailUUID = tailNode.Next
			}

			// Link new node to the real tail.
			tailNode.Next = newNodeUUID

			tail_Node_Bytes, err := json.Marshal(tailNode)
			if err != nil {
				return err
			}
			tailCT := userlib.SymEnc(fileEncKey[:16], userlib.RandomBytes(16), tail_Node_Bytes)
			tailTag, err := userlib.HMACEval(fileMacKey[:16], tailCT)
			if err != nil {
				return err
			}
			newTailBlob, err := json.Marshal(struct {
				C   []byte
				Tag []byte
			}{
				C:   tailCT,
				Tag: tailTag,
			})
			if err != nil {
				return err
			}
			userlib.DatastoreSet(tailUUID, newTailBlob)

			// Now meta.Current_Node really is the tail.
			meta.Current_Node = newNodeUUID
		}

		// 1c. Re-encrypt and store updated FileMeta.
		metaBytes, err := json.Marshal(meta)
		if err != nil {
			return err
		}
		metaCT := userlib.SymEnc(metaEncKey, userlib.RandomBytes(16), metaBytes)
		metaTag, err := userlib.HMACEval(metaMacKey, metaCT)
		if err != nil {
			return err
		}
		finalMeta, err := json.Marshal(struct {
			C   []byte
			Tag []byte
		}{
			C:   metaCT,
			Tag: metaTag,
		})
		if err != nil {
			return err
		}
		userlib.DatastoreSet(fileMetaUUID, finalMeta)

		return nil
	}

	// NON-OWNER path via handle
	// No FileMeta found under this user's root-derived UUID; treat as non-owner.
	handleKey, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("handle-uuid"+filename))
	if err != nil {
		return err
	}
	handleUUID, err := uuid.FromBytes(handleKey[:16])
	if err != nil {
		return err
	}

	handleBlob, ok := userlib.DatastoreGet(handleUUID)
	if !ok {
		return errors.New("file not found for this user")
	}

	// Decrypt & verify the handle (SharedFile).
	var handleWrapped struct {
		C   []byte
		Tag []byte
	}
	if err := json.Unmarshal(handleBlob, &handleWrapped); err != nil {
		return err
	}

	hEncMat, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("handle-enc"+filename))
	if err != nil {
		return err
	}
	hMacMat, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("handle-mac"+filename))
	if err != nil {
		return err
	}
	hEncKey := hEncMat[:16]
	hMacKey := hMacMat[:16]

	hTagCheck, err := userlib.HMACEval(hMacKey, handleWrapped.C)
	if err != nil {
		return err
	}
	if !userlib.HMACEqual(hTagCheck, handleWrapped.Tag) {
		return errors.New("shared file handle tampered")
	}

	handlePlain := userlib.SymDec(hEncKey, handleWrapped.C)
	var handle SharedFile
	if err := json.Unmarshal(handlePlain, &handle); err != nil {
		return err
	}

	// Fetch and verify the Invitation this handle points to.
	invBlobBytes, ok := userlib.DatastoreGet(handle.InviteUUID)
	if !ok {
		return errors.New("invitation not found")
	}
	// we need an invite blob so that the person doing all the invite thingys can decrpyt and doo all the magical shito
	type invitationBlob struct {
		EncData   []byte
		Signature []byte
		Sender    string
	}

	var blob invitationBlob
	if err := json.Unmarshal(invBlobBytes, &blob); err != nil {
		return err
	}

	// Verify signature using the sender's verify key (same convention as AcceptInvitation).
	senderVerifyKey, ok := userlib.KeystoreGet(blob.Sender + "_dsv")
	if !ok {
		return errors.New("sender verify key not found")
	}
	if err := userlib.DSVerify(senderVerifyKey, blob.EncData, blob.Signature); err != nil {
		return errors.New("invalid invitation signature")
	}

	// Decrypt inner Invitation with this user's private PKE key.
	invPlain, err := userlib.PKEDec(userdata.PrivEncKey, blob.EncData)
	if err != nil {
		return errors.New("could not decrypt invitation")
	}

	inv, err := unpackInvitation(invPlain)
	if err != nil {
		return err
	}

	// Respect revocation flag.
	if inv.Revoked {
		return errors.New("access to this file has been revoked")
	}

	// Now we have per-file keys and the base node address from the Invitation.
	fileEncKey := inv.EncKey
	fileMacKey := inv.MacKey
	curUUID := inv.FileMetaData // should be the base FileNode for this file

	if curUUID == uuid.Nil {
		return errors.New("invitation points to nil node")
	}

	// Walk the linked list to find the current tail node.
	tailUUID := curUUID
	var tailNode FileNode

	for {
		nodeBlob, ok := userlib.DatastoreGet(tailUUID)
		if !ok {
			return errors.New("file node missing from datastore")
		}

		var nodeWrapped struct {
			C   []byte
			Tag []byte
		}
		if err := json.Unmarshal(nodeBlob, &nodeWrapped); err != nil {
			return err
		}

		nodeTagCheck, err := userlib.HMACEval(fileMacKey[:16], nodeWrapped.C)
		if err != nil {
			return err
		}
		if !userlib.HMACEqual(nodeTagCheck, nodeWrapped.Tag) {
			return errors.New("file node tampered")
		}

		nodePlain := userlib.SymDec(fileEncKey[:16], nodeWrapped.C)
		var node FileNode
		if err := json.Unmarshal(nodePlain, &node); err != nil {
			return err
		}

		if node.Next == uuid.Nil {
			// Found tail.
			tailNode = node
			break
		}
		// Otherwise keep walking.
		tailUUID = node.Next
	}

	// Create the new appended node.
	newNodeUUID := uuid.New()
	newNode := FileNode{
		Content: content,
		Next:    uuid.Nil,
	}
	newNodeBytes, err := json.Marshal(newNode)
	if err != nil {
		return err
	}
	newCT := userlib.SymEnc(fileEncKey[:16], userlib.RandomBytes(16), newNodeBytes)
	newTag, err := userlib.HMACEval(fileMacKey[:16], newCT)
	if err != nil {
		return err
	}
	newNodeBlob, err := json.Marshal(struct {
		C   []byte
		Tag []byte
	}{
		C:   newCT,
		Tag: newTag,
	})
	if err != nil {
		return err
	}
	userlib.DatastoreSet(newNodeUUID, newNodeBlob)

	// Update old tail's Next to point to the new node.
	tailNode.Next = newNodeUUID
	tailNodeBytes, err := json.Marshal(tailNode)
	if err != nil {
		return err
	}
	tailCT := userlib.SymEnc(fileEncKey[:16], userlib.RandomBytes(16), tailNodeBytes)
	tailTag, err := userlib.HMACEval(fileMacKey[:16], tailCT)
	if err != nil {
		return err
	}
	updatedTailBlob, err := json.Marshal(struct {
		C   []byte
		Tag []byte
	}{
		C:   tailCT,
		Tag: tailTag,
	})
	if err != nil {
		return err
	}
	userlib.DatastoreSet(tailUUID, updatedTailBlob)

	return nil
}

func (userdata *User) LoadFile(filename string) (content []byte, err error) {
	// if filename == "" {
	// 	return nil, errors.New("empty filename")
	// }

	// =======================
	// 1. Try OWNER path first
	// =======================
	fileKey, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("file-uuid"+filename))
	if err != nil {
		return nil, err
	}
	fileMetaUUID, err := uuid.FromBytes(fileKey[:16])
	if err != nil {
		return nil, err
	}

	metaEncMaterial, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("file-meta-enc"+filename))
	if err != nil {
		return nil, err
	}
	metaMacMaterial, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("file-meta-mac"+filename))
	if err != nil {
		return nil, err
	}
	metaEncKey := metaEncMaterial[:16]
	metaMacKey := metaMacMaterial[:16]

	metaBlob, ok := userlib.DatastoreGet(fileMetaUUID)
	if ok {
		// We might be the owner of this file.
		var wrapped struct {
			C   []byte
			Tag []byte
		}
		if err := json.Unmarshal(metaBlob, &wrapped); err != nil {
			return nil, err
		}

		tagCheck, err := userlib.HMACEval(metaMacKey, wrapped.C)
		if err != nil {
			return nil, err
		}
		if !userlib.HMACEqual(tagCheck, wrapped.Tag) {
			return nil, errors.New("file metadata tampered")
		}

		metaPlain := userlib.SymDec(metaEncKey, wrapped.C)

		var meta FileMeta
		if err := json.Unmarshal(metaPlain, &meta); err != nil {
			return nil, err
		}

		// Owner path: read from meta.Base_Node using per-file keys.
		baseUUID := meta.Base_Node
		fileEncKey := meta.EncKey
		fileMacKey := meta.MACKey

		// Empty file? Return empty slice.
		if baseUUID == uuid.Nil {
			return []byte{}, nil
		}

		// Traverse linked list from base node, concatenating Content.
		var result []byte
		curUUID := baseUUID

		for curUUID != uuid.Nil {
			nodeBlob, ok := userlib.DatastoreGet(curUUID)
			if !ok {
				return nil, errors.New("file node missing from datastore")
			}

			var nodeWrapped struct {
				C   []byte
				Tag []byte
			}
			if err := json.Unmarshal(nodeBlob, &nodeWrapped); err != nil {
				return nil, err
			}

			nodeTagCheck, err := userlib.HMACEval(fileMacKey[:16], nodeWrapped.C)
			if err != nil {
				return nil, err
			}
			if !userlib.HMACEqual(nodeTagCheck, nodeWrapped.Tag) {
				return nil, errors.New("file node tampered")
			}

			nodePlain := userlib.SymDec(fileEncKey[:16], nodeWrapped.C)
			var node FileNode
			if err := json.Unmarshal(nodePlain, &node); err != nil {
				return nil, err
			}

			result = append(result, node.Content...)
			curUUID = node.Next
		}

		return result, nil
	}

	// ===========================
	// 2. NON-OWNER path via handle
	// ===========================
	handleKey, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("handle-uuid"+filename))
	if err != nil {
		return nil, err
	}
	handleUUID, err := uuid.FromBytes(handleKey[:16])
	if err != nil {
		return nil, err
	}

	handleBlob, ok := userlib.DatastoreGet(handleUUID)
	if !ok {
		return nil, errors.New("file not found for this user")
	}

	// Decrypt & verify the handle (SharedFile).
	var handleWrapped struct {
		C   []byte
		Tag []byte
	}
	if err := json.Unmarshal(handleBlob, &handleWrapped); err != nil {
		return nil, err
	}

	hEncMat, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("handle-enc"+filename))
	if err != nil {
		return nil, err
	}
	hMacMat, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("handle-mac"+filename))
	if err != nil {
		return nil, err
	}
	hEncKey := hEncMat[:16]
	hMacKey := hMacMat[:16]

	hTagCheck, err := userlib.HMACEval(hMacKey, handleWrapped.C)
	if err != nil {
		return nil, err
	}
	if !userlib.HMACEqual(hTagCheck, handleWrapped.Tag) {
		return nil, errors.New("shared file handle tampered")
	}

	handlePlain := userlib.SymDec(hEncKey, handleWrapped.C)
	var handle SharedFile
	if err := json.Unmarshal(handlePlain, &handle); err != nil {
		return nil, err
	}

	// Fetch and verify the Invitation this handle points to.
	invBlobBytes, ok := userlib.DatastoreGet(handle.InviteUUID)
	if !ok {
		return nil, errors.New("invitation not found")
	}

	type invitationBlob struct {
		EncData   []byte
		Signature []byte
		Sender    string
	}

	var blob invitationBlob
	if err := json.Unmarshal(invBlobBytes, &blob); err != nil {
		return nil, err
	}

	// Verify signature using sender's public DS verify key.
	senderVerifyKey, ok := userlib.KeystoreGet(blob.Sender + "_dsv")
	if !ok {
		return nil, errors.New("sender verify key not found")
	}
	if err := userlib.DSVerify(senderVerifyKey, blob.EncData, blob.Signature); err != nil {
		return nil, errors.New("invalid invitation signature")
	}

	// Decrypt inner Invitation with this user's private key.
	invPlain, err := userlib.PKEDec(userdata.PrivEncKey, blob.EncData)
	if err != nil {
		return nil, errors.New("could not decrypt invitation")
	}

	inv, err := unpackInvitation(invPlain)
	if err != nil {
		return nil, err
	}

	// Respect revocation.
	if inv.Revoked {
		return nil, errors.New("access to this file has been revoked")
	}

	// Non-owner path: read from inv.NodeAddress using per-file keys.
	baseUUID := inv.FileMetaData
	fileEncKey := inv.EncKey
	fileMacKey := inv.MacKey

	if baseUUID == uuid.Nil {
		return nil, errors.New("invitation points to nil base node")
	}

	var result []byte
	curUUID := baseUUID

	for curUUID != uuid.Nil {
		nodeBlob, ok := userlib.DatastoreGet(curUUID)
		if !ok {
			return nil, errors.New("file node missing from datastore")
		}

		var nodeWrapped struct {
			C   []byte
			Tag []byte
		}
		if err := json.Unmarshal(nodeBlob, &nodeWrapped); err != nil {
			return nil, err
		}

		nodeTagCheck, err := userlib.HMACEval(fileMacKey[:16], nodeWrapped.C)
		if err != nil {
			return nil, err
		}
		if !userlib.HMACEqual(nodeTagCheck, nodeWrapped.Tag) {
			return nil, errors.New("file node tampered")
		}

		nodePlain := userlib.SymDec(fileEncKey[:16], nodeWrapped.C)
		var node FileNode
		if err := json.Unmarshal(nodePlain, &node); err != nil {
			return nil, err
		}

		result = append(result, node.Content...)
		curUUID = node.Next
	}

	return result, nil
}

//----------------------------------create invitation------------------------------------------------------------------

func (userdata *User) CreateInvitation(filename string, recipientUsername string) (
	invitationPtr uuid.UUID, err error) {
	// if filename == "" {
	// 	return uuid.Nil, errors.New("filename required")
	// }
	if recipientUsername == "" {
		return uuid.Nil, errors.New("recipient username required")
	}
	if recipientUsername == userdata.Username {
		return uuid.Nil, errors.New("cannot invite yourself")
	}

	// Get recipient's public key up front (used in both owner + shared paths).
	recipientPubKey, ok := userlib.KeystoreGet(recipientUsername + "_pke")
	if !ok {
		return uuid.Nil, errors.New("recipient public key not found")
	}

	// Helper: sign+store an Invitation and return its UUID.
	makeAndStoreInvitation := func(ownerName string, metaUUID uuid.UUID, fileEncKey, fileMacKey []byte) (uuid.UUID, error) {
		inv := Invitation{
			OGOwner:      ownerName,
			FileMetaData: metaUUID,
			MacKey:       fileMacKey,
			EncKey:       fileEncKey,
			Revoked:      false,
		}

		invBytes, err := packInvitation(inv)
		if err != nil {
			return uuid.Nil, err
		}

		encInvite, err := userlib.PKEEnc(recipientPubKey, invBytes)
		if err != nil {
			return uuid.Nil, err
		}

		// if err != nil {
		// 	return uuid.Nil, err
		// }

		signature, err := userlib.DSSign(userdata.PrivSignatureKey, encInvite)
		if err != nil {
			return uuid.Nil, err
		}

		type invitationBlob struct {
			EncData   []byte
			Signature []byte
			Sender    string
		}

		blob := invitationBlob{
			EncData:   encInvite,
			Signature: signature,
			Sender:    userdata.Username,
		}

		blobBytes, err := json.Marshal(blob)
		if err != nil {
			return uuid.Nil, err
		}

		newPtr := uuid.New()
		userlib.DatastoreSet(newPtr, blobBytes)
		return newPtr, nil
	}

	//  OWNER PATH: try to load FileMeta using this user's UserRootkey.

	fileKey, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("file-uuid"+filename))
	if err != nil {
		return uuid.Nil, err
	}
	fileMetaUUID, err := uuid.FromBytes(fileKey[:16])
	if err != nil {
		return uuid.Nil, err
	}

	metaBlob, ok := userlib.DatastoreGet(fileMetaUUID)
	if ok {
		// We *might* be the owner. Try to decrypt FileMeta with owner metadata keys.
		metaEncMat, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("file-meta-enc"+filename))
		if err != nil {
			return uuid.Nil, err
		}
		metaMacMat, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("file-meta-mac"+filename))
		if err != nil {
			return uuid.Nil, err
		}
		metaEncKey := metaEncMat[:16]
		metaMacKey := metaMacMat[:16]

		var wrapped struct {
			C   []byte
			Tag []byte
		}
		if err := json.Unmarshal(metaBlob, &wrapped); err != nil {
			return uuid.Nil, err
		}

		tagCheck, err := userlib.HMACEval(metaMacKey, wrapped.C)
		if err != nil {
			return uuid.Nil, err
		}
		if !userlib.HMACEqual(tagCheck, wrapped.Tag) {
			// We found something at the expected location but MAC failed → integrity error,
			// not a "just fall back" case.
			return uuid.Nil, errors.New("file metadata tampered")
		}

		metaPlain := userlib.SymDec(metaEncKey, wrapped.C)

		var meta FileMeta
		if err := json.Unmarshal(metaPlain, &meta); err != nil {
			return uuid.Nil, err
		}

		// At this point we are the original owner.
		// Create a new Invitation using the per-file keys from FileMeta.
		invitationPtr, err = makeAndStoreInvitation(meta.Owner, meta.Base_Node, meta.EncKey, meta.MACKey)
		if err != nil {
			return uuid.Nil, err
		}

		//  Update Owner_map in FileMeta so RevokeAccess can later see direct sharees. ---
		//  before loading owner map
		// if meta.Owner_map == uuid.Nil {
		// 		meta.Owner_map = uuid.New()
		// 		// re-store updated FileMeta with the new Owner_map pointer
		// 		metaBytes, err := json.Marshal(meta)
		// 		if err != nil {
		// 			return uuid., err
		// 		}
		// 		metaCT := userlib.SymEnc(mEncKey, usib.RandomBytes(16), metaBytes)
		// 		metaTag, err := usrlib.HMACEval(metaMcKey, metaCT)
		// 		if err != nil {
		// 			return uui.Nil, err
		// 		}
		// you are more that I need you are more than just a dream
		if meta.Owner_map == uuid.Nil {
			meta.Owner_map = uuid.New()
			metaBytes, err := json.Marshal(meta)
			if err != nil {
				return uuid.Nil, err
			}
			metaCT := userlib.SymEnc(metaEncKey, userlib.RandomBytes(16), metaBytes)
			metaTag, err := userlib.HMACEval(metaMacKey, metaCT)
			if err != nil {
				return uuid.Nil, err
			}
			newMetaBlob, err := json.Marshal(struct {
				C   []byte
				Tag []byte
			}{
				C:   metaCT,
				Tag: metaTag,
			})
			if err != nil {
				return uuid.Nil, err
			}
			userlib.DatastoreSet(fileMetaUUID, newMetaBlob)

		}

		// we Use a single invitation pointer per username so wedont get gucked later on
		ownerMap := make(map[string]uuid.UUID)
		if omBlob, ok := userlib.DatastoreGet(meta.Owner_map); ok {
			var owWrapped struct {
				C   []byte
				Tag []byte
			}
			if err := json.Unmarshal(omBlob, &owWrapped); err != nil {
				return uuid.Nil, err
			}

			owTagCheck, err := userlib.HMACEval(meta.MACKey[:16], owWrapped.C)
			if err != nil {
				return uuid.Nil, err
			}
			if !userlib.HMACEqual(owTagCheck, owWrapped.Tag) {
				return uuid.Nil, errors.New("owner map tampered")
			}

			owPlain := userlib.SymDec(meta.EncKey[:16], owWrapped.C)
			if err := json.Unmarshal(owPlain, &ownerMap); err != nil {
				return uuid.Nil, err
			}
		}

		// Store/overwrite the invitation pointer for this recipient.
		ownerMap[recipientUsername] = invitationPtr

		// Re-encrypt and store updated owner map using per-file keys.
		omPlain, err := json.Marshal(ownerMap)
		if err != nil {
			return uuid.Nil, err
		}
		omCT := userlib.SymEnc(meta.EncKey[:16], userlib.RandomBytes(16), omPlain)
		omTag, err := userlib.HMACEval(meta.MACKey[:16], omCT)
		if err != nil {
			return uuid.Nil, err
		}
		finalOM, err := json.Marshal(struct {
			C   []byte
			Tag []byte
		}{
			C:   omCT,
			Tag: omTag,
		})
		if err != nil {
			return uuid.Nil, err
		}

		userlib.DatastoreSet(meta.Owner_map, finalOM)

		return invitationPtr, nil
	}

	//
	// 2) SHARED PATH: we are *not* the owner, but we might have a handle.

	// Find our handle for this filename.
	// remeber accorifng to the spec a file name can be "" so I hope i dont get fucked for that over the spec
	//later on in life
	handleKey, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("handle-uuid"+filename))
	if err != nil {
		return uuid.Nil, err
	}
	handleUUID, err := uuid.FromBytes(handleKey[:16])
	if err != nil {
		return uuid.Nil, err
	}

	handleBlob, ok := userlib.DatastoreGet(handleUUID)
	if !ok {
		return uuid.Nil, errors.New("chat me might be cooked, so such filename")
	}

	// Derive handle enc/mac keys.
	hEncMat, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("handle-enc"+filename))
	if err != nil {
		return uuid.Nil, err
	}
	hMacMat, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("handle-mac"+filename))
	if err != nil {
		return uuid.Nil, err
	}
	hEncKey := hEncMat[:16]
	hMacKey := hMacMat[:16]

	var hWrapped struct {
		C   []byte
		Tag []byte
	}
	if err := json.Unmarshal(handleBlob, &hWrapped); err != nil {
		return uuid.Nil, err
	}

	h_Tag_Check, err := userlib.HMACEval(hMacKey, hWrapped.C)
	if err != nil {
		return uuid.Nil, err
	}
	if !userlib.HMACEqual(h_Tag_Check, hWrapped.Tag) {
		return uuid.Nil, errors.New("someone_is_touching the handle")
	}

	hP_lain := userlib.SymDec(hEncKey, hWrapped.C)

	var handle SharedFile
	if err := json.Unmarshal(hP_lain, &handle); err != nil {
		return uuid.Nil, err
	}

	// Follow the existing invitation we were given.
	invBlobBytes, ok := userlib.DatastoreGet(handle.InviteUUID)
	if !ok {
		return uuid.Nil, errors.New("original invitation not found")
	}

	type invitationBlob struct {
		EncData   []byte
		Signature []byte
		Sender    string
	}

	var inBlob invitationBlob
	if err := json.Unmarshal(invBlobBytes, &inBlob); err != nil {
		return uuid.Nil, err
	}

	// Verify the sender's signature on the existing invitation.
	senderVerifyKey, ok := userlib.KeystoreGet(inBlob.Sender + "_dsv")
	if !ok {
		return uuid.Nil, errors.New("sender verify key not found")
	}
	if err := userlib.DSVerify(senderVerifyKey, inBlob.EncData, inBlob.Signature); err != nil {
		return uuid.Nil, errors.New("invitation signature invalid")
	}

	// Decrypt the inner Invitation with *our* private key.
	invPlain, err := userlib.PKEDec(userdata.PrivEncKey, inBlob.EncData)
	if err != nil {
		return uuid.Nil, err
	}
	inv, err := unpackInvitation(invPlain)
	if err != nil {
		return uuid.Nil, err
	}

	if inv.Revoked {
		return uuid.Nil, errors.New("invitation already revoked")
	}

	// Build a new invitation for the new recipient using the same per-file keys.
	return makeAndStoreInvitation(inv.OGOwner, inv.FileMetaData, inv.EncKey, inv.MacKey)

}

func (userdata *User) AcceptInvitation(senderUsername string, invitationPtr uuid.UUID, filename string) error {
	if filename == "" {
		return errors.New("empty filename")
	}
	if senderUsername == "" {
		return errors.New("empty sender username")
	}

	fileKey, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("file-uuid"+filename))
	if err != nil {
		return err
	}

	fileMetaUUID, err := uuid.FromBytes(fileKey[:16])
	if err != nil {
		return err
	}

	if _, ok := userlib.DatastoreGet(fileMetaUUID); ok {
		return errors.New("file name already exists")
	}

	handleKey, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("handle-uuid"+filename))
	if err != nil {
		return err
	}
	handleUUID, err := uuid.FromBytes(handleKey[:16])
	if err != nil {
		return err
	}

	if _, ok := userlib.DatastoreGet(handleUUID); ok {
		return errors.New("a file with this name already exists in user's namespace")
	}

	invBlobBytes, ok := userlib.DatastoreGet(invitationPtr)
	if !ok {
		return errors.New("invitation not found")
	}

	type invitationBlob struct {
		EncData   []byte
		Signature []byte
		Sender    string
	}

	var blob invitationBlob
	if err := json.Unmarshal(invBlobBytes, &blob); err != nil {
		return err
	}

	if blob.Sender != senderUsername {
		return errors.New("sender username mismatch")
	}

	sender_Verify_Key, ok := userlib.KeystoreGet(senderUsername + "_dsv")
	if !ok {
		return errors.New("sender verify key not found")
	}

	if err := userlib.DSVerify(sender_Verify_Key, blob.EncData, blob.Signature); err != nil {
		return errors.New("invalid invitation signature")
	}

	invPlain, err := userlib.PKEDec(userdata.PrivEncKey, blob.EncData)
	if err != nil {
		return errors.New("could not decrypt invitation")
	}

	inv, err := unpackInvitation(invPlain)
	if err != nil {
		return err
	}

	// If the invitation has been marked revoked, don't accept it.
	if inv.Revoked {
		return errors.New("invitation has been revoked")
	}

	handle := SharedFile{
		OGOwner:     inv.OGOwner,
		NodeAddress: inv.FileMetaData, // MetaUUID of FileMeta
		InviteUUID:  invitationPtr,    // live capability we’ll follow for keys
	}

	handleBytes, err := json.Marshal(handle)
	if err != nil {
		return err
	}

	hEncMat, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("handle-enc"+filename))
	if err != nil {
		return err
	}
	hMacMat, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("handle-mac"+filename))
	if err != nil {
		return err
	}
	hEncKey := hEncMat[:16]
	hMacKey := hMacMat[:16]

	h_C_T := userlib.SymEnc(hEncKey, userlib.RandomBytes(16), handleBytes)
	hTag, err := userlib.HMACEval(hMacKey, h_C_T)
	if err != nil {
		return err
	}

	finalHandleBlob, err := json.Marshal(struct {
		C   []byte
		Tag []byte
	}{
		C:   h_C_T,
		Tag: hTag,
	})
	if err != nil {
		return err
	}

	//mew moew moew moew moew moew moemw o mew mowm woemw eowm o
	// /
	//
	// I hate htis projet makes me owanna kill ymself yeah
	// I think i died and came back to life fot htis pen
	//
	// ore the handle at the deterministic handleUUID.
	userlib.DatastoreSet(handleUUID, finalHandleBlob)

	return nil

}

//return nil
//pookie gpt was ued to help with the code for revoke (not the design) the design was all your boi and lowkey am proud of
// doing so, also pookie gpt has been used throughout the project for syntax help and stuff since we dont really teach GO
//

func (userdata *User) RevokeAccess(filename string, recipientUsername string) error {

	if filename == "" || recipientUsername == "" {
		return errors.New("filename and recipient username must be non-empty")
	}

	// Same bull poo poo  as StoreFile  AppendToFile  LoadFile.
	fileKey, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("file-uuid"+filename))
	if err != nil {
		return err
	}
	fileMetaUUID, err := uuid.FromBytes(fileKey[:16])
	if err != nil {
		return err
	}

	metaEncMaterial, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("file-meta-enc"+filename))
	if err != nil {
		return err
	}
	metaMacMaterial, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("file-meta-mac"+filename))
	if err != nil {
		return err
	}
	metaEncKey := metaEncMaterial[:16]
	meta_Mac_Key := metaMacMaterial[:16]

	metaBlob, ok := userlib.DatastoreGet(fileMetaUUID)
	if !ok {
		return errors.New("file not found")
	}

	var metaWrapped struct {
		C   []byte
		Tag []byte
	}
	if err := json.Unmarshal(metaBlob, &metaWrapped); err != nil {
		return err
	}

	meta_Tag_Check, err := userlib.HMACEval(meta_Mac_Key, metaWrapped.C)
	if err != nil {
		return err
	}
	if !userlib.HMACEqual(meta_Tag_Check, metaWrapped.Tag) {
		return errors.New("file metadata tampered")
	}

	metaPlain := userlib.SymDec(metaEncKey, metaWrapped.C)

	var meta FileMeta
	if err := json.Unmarshal(metaPlain, &meta); err != nil {
		return err
	}

	// Only the original owner is allowed to revoke.
	if meta.Owner != userdata.Username {
		return errors.New("only the original owner may revoke access")
	}

	if meta.Owner_map == uuid.Nil {
		return errors.New("no shares for this file")
	}

	//  Load owner_map

	ownerMapUUID := meta.Owner_map

	ownerMapBlob, ok := userlib.DatastoreGet(ownerMapUUID)
	if !ok {
		return errors.New("owner map missing")
	}

	var ownerMapWrapped struct {
		C   []byte
		Tag []byte
	}
	if err := json.Unmarshal(ownerMapBlob, &ownerMapWrapped); err != nil {
		return err
	}

	// Derive enc/mac keys for the owner map from the user root key + filename.
	// ownerMapEncMat, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("owner-map-enc"+filename))
	// if err != nil {
	// 	return err
	// }
	// ownerMapMacMat, err := userlib.HashKDF(userdata.UserRootkey[:16], []byte("owner-map-mac"+filename))
	// if err != nil {
	// 	return err
	// }
	// ownerMapEncKey := ownerMapEncMat[:16]
	// ownerMapMacKey := ownerMapMacMat[:16]

	ownerMapTagCheck, err := userlib.HMACEval(meta.MACKey[:16], ownerMapWrapped.C)
	if err != nil {
		return err
	}
	if !userlib.HMACEqual(ownerMapTagCheck, ownerMapWrapped.Tag) {
		return errors.New("owner map tampered")
	}

	ownerMapPlain := userlib.SymDec(meta.EncKey[:16], ownerMapWrapped.C)

	var ownerMap map[string]uuid.UUID
	if err := json.Unmarshal(ownerMapPlain, &ownerMap); err != nil {
		return err
	}

	revokedPtr, ok := ownerMap[recipientUsername]
	if !ok {
		return errors.New("recipient was not directly shared this file by owner")
	}

	remainingUsers := make([]string, 0)
	for uname := range ownerMap {
		if uname != recipientUsername {
			remainingUsers = append(remainingUsers, uname)
		}
	}

	oldEncKey := meta.EncKey
	oldMacKey := meta.MACKey

	if len(oldEncKey) < 16 || len(oldMacKey) < 16 {
		return errors.New("invalid old file keys")
	}

	newEncKey := userlib.RandomBytes(16)
	newMacKey := userlib.RandomBytes(16)

	curUUID := meta.Base_Node
	for curUUID != uuid.Nil {
		nodeBlob, ok := userlib.DatastoreGet(curUUID)
		if !ok {
			return errors.New("file node missing during rekey")
		}

		var nodeWrapped struct {
			C   []byte
			Tag []byte
		}
		if err := json.Unmarshal(nodeBlob, &nodeWrapped); err != nil {
			return err
		}

		nodeTagCheck, err := userlib.HMACEval(oldMacKey[:16], nodeWrapped.C)
		if err != nil {
			return err
		}
		if !userlib.HMACEqual(nodeTagCheck, nodeWrapped.Tag) {
			return errors.New("file node tampered during rekey")
		}

		nodePlain := userlib.SymDec(oldEncKey[:16], nodeWrapped.C)

		var node FileNode
		if err := json.Unmarshal(nodePlain, &node); err != nil {
			return err
		}

		// Re-encrypt under new keys.
		nodeBytes, err := json.Marshal(node)
		if err != nil {
			return err
		}

		nodeCT := userlib.SymEnc(newEncKey[:16], userlib.RandomBytes(16), nodeBytes)
		nodeTag, err := userlib.HMACEval(newMacKey[:16], nodeCT)
		if err != nil {
			return err
		}

		newNodeBlob, err := json.Marshal(struct {
			C   []byte
			Tag []byte
		}{
			C:   nodeCT,
			Tag: nodeTag,
		})
		if err != nil {
			return err
		}

		userlib.DatastoreSet(curUUID, newNodeBlob)

		// Continue traversal using the Next pointer.
		curUUID = node.Next
	}

	// Update FileMeta with new keys and write it back.
	meta.EncKey = newEncKey
	meta.MACKey = newMacKey

	metaBytes, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	metaCT := userlib.SymEnc(metaEncKey[:16], userlib.RandomBytes(16), metaBytes)
	metaNewTag, err := userlib.HMACEval(meta_Mac_Key[:16], metaCT)
	if err != nil {
		return err
	}
	newMetaBlob, err := json.Marshal(struct {
		C   []byte
		Tag []byte
	}{
		C:   metaCT,
		Tag: metaNewTag,
	})
	if err != nil {
		return err
	}
	userlib.DatastoreSet(fileMetaUUID, newMetaBlob)

	for _, child := range remainingUsers {
		invPtr := ownerMap[child]

		recipientPubKey, ok := userlib.KeystoreGet(child + "_pke")
		if !ok {
			// If they somehow don't exist anymore, just skip them.
			continue
		}

		// Build fresh Invitation with new keys.
		invStruct := Invitation{
			OGOwner:      meta.Owner,
			FileMetaData: meta.Base_Node,
			MacKey:       newMacKey,
			EncKey:       newEncKey,
			Revoked:      false,
		}

		packed, err := packInvitation(invStruct)
		if err != nil {
			return err
		}

		encInvite, err := userlib.PKEEnc(recipientPubKey, packed)
		if err != nil {
			return err
		}

		sig, err := userlib.DSSign(userdata.PrivSignatureKey, encInvite)
		if err != nil {
			return err
		}

		// Same outer wrapper type used elsewhere.
		outer := struct {
			EncData   []byte
			Signature []byte
			Sender    string
		}{
			EncData:   encInvite,
			Signature: sig,
			Sender:    meta.Owner,
		}

		outerBytes, err := json.Marshal(outer)
		if err != nil {
			return err
		}

		userlib.DatastoreSet(invPtr, outerBytes)
	}

	//

	// Not strictly required, since they only have old keys and we've re-encrypted all nodes.
	// now here is the sexy part because of the revoked boolean of our inv struct
	// we can quickly check for revocation and then if it is revoked then we just return in load and append
	if revokedPubKey, ok := userlib.KeystoreGet(recipientUsername + "_pke"); ok {
		zeroKey := make([]byte, 16)

		revokedInv := Invitation{
			OGOwner:      meta.Owner,
			FileMetaData: meta.Base_Node,
			MacKey:       zeroKey,
			EncKey:       zeroKey,
			Revoked:      true,
		}

		packed, err := packInvitation(revokedInv)
		if err == nil {
			encInvite, err := userlib.PKEEnc(revokedPubKey, packed)
			if err == nil {
				sig, err := userlib.DSSign(userdata.PrivSignatureKey, encInvite)
				if err == nil {
					outer := struct {
						EncData   []byte
						Signature []byte
						Sender    string
					}{
						EncData:   encInvite,
						Signature: sig,
						Sender:    meta.Owner,
					}
					outerBytes, err := json.Marshal(outer)
					if err == nil {
						userlib.DatastoreSet(revokedPtr, outerBytes)
					}
				}
			}
		}
	}

	delete(ownerMap, recipientUsername)

	ownerMapBytes, err := json.Marshal(ownerMap)
	if err != nil {
		return err
	}
	ownerMapCT := userlib.SymEnc(meta.EncKey[:16], userlib.RandomBytes(16), ownerMapBytes)
	ownerMapTag, err := userlib.HMACEval(meta.MACKey[:16], ownerMapCT)
	if err != nil {
		return err
	}
	newOwnerMapBlob, err := json.Marshal(struct {
		C   []byte
		Tag []byte
	}{
		C:   ownerMapCT,
		Tag: ownerMapTag,
	})
	if err != nil {
		return err
	}
	userlib.DatastoreSet(ownerMapUUID, newOwnerMapBlob)

	return nil
}

// GOOD shit boyz GG
// Bravo 6 going dark....
