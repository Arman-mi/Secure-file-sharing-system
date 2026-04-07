# Secure File Sharing System (CS 161 Project 2)

This repository contains my implementation of a secure file sharing system. The project builds a client-side encrypted file store on top of an untrusted datastore and a trusted keystore for public keys. Users can store, append to, load, share, and revoke access to files while preserving confidentiality and integrity against a malicious datastore.

Spec link: [Project 2 spec](https://fa25.cs161.org/proj2/)

**Quick Summary**
1. Client-side encryption and integrity for all file and user data.
2. Efficient append via an encrypted linked-list of file nodes.
3. Share and revoke access using signed, encrypted invitations.

**How It Works**
1. **User initialization and login**
   - `InitUser` derives a base key from the password using Argon2 with a random salt, then derives separate encryption, MAC, and root keys.
   - The user struct is serialized, encrypted, MACed, and stored at a deterministic UUID derived from the username.
   - Public keys are stored in the keystore under names like `username_pke` and `username_dsv`.
2. **File storage**
   - Each file has a `FileMeta` record stored at a deterministic UUID derived from the user’s root key and the filename.
   - `FileMeta` contains per-file symmetric keys plus pointers to the head/tail of a linked list of `FileNode` objects.
   - Each `FileNode` stores encrypted content and a pointer to the next node, allowing O(1) append.
3. **Append**
   - `AppendToFile` creates a new encrypted `FileNode`, updates the previous tail node’s `Next` pointer, and updates `FileMeta`.
4. **Load**
   - Owners decrypt `FileMeta` and traverse the linked list, verifying a MAC on each node.
   - Shared users decrypt a local “handle,” verify the invitation signature, decrypt the invitation to obtain per-file keys, then traverse nodes.
5. **Sharing**
   - `CreateInvitation` packages per-file keys and a base-node pointer into an `Invitation`, encrypts it to the recipient, and signs it.
   - The recipient stores a local `SharedFile` handle (encrypted/MACed under their root key) that points to the invitation UUID.
6. **Revocation**
   - Only the original owner can revoke.
   - `RevokeAccess` rekeys the file by decrypting each node with old keys and re-encrypting with fresh keys.
   - Invitations for remaining users are reissued with new keys; the revoked user’s invitation is overwritten with a revoked marker.

**Public API (package `client`)**
1. `InitUser(username, password)`
2. `GetUser(username, password)`
3. `StoreFile(filename, content)`
4. `LoadFile(filename)`
5. `AppendToFile(filename, content)`
6. `CreateInvitation(filename, recipientUsername)`
7. `AcceptInvitation(senderUsername, invitationPtr, filename)`
8. `RevokeAccess(filename, recipientUsername)`

**Repository Layout**
1. `client/client.go` – implementation
2. `client/client_unittest.go` – optional white-box unit tests
3. `client_test/client_test.go` – black-box integration tests
4. `my_tests/my_test.go` – additional tests

**Running Tests**
```bash
go test ./...
```

**Notes**
1. This is a course project and is not production hardened.
2. The system relies on `github.com/cs161-staff/project2-userlib` for cryptographic primitives and the datastore/keystore API.
