package iam

import (
	"context"
	"testing"

	"github.com/hashicorp/watchtower/internal/db"
	"google.golang.org/protobuf/proto"
	"gotest.tools/assert"
)

func Test_NewAuthMethod(t *testing.T) {
	db.StartTest()
	t.Parallel()
	cleanup, url := db.SetupTest(t, "../db/migrations/postgres")
	defer cleanup()
	defer db.CompleteTest() // must come after the "defer cleanup()"
	conn, err := db.TestConnection(url)
	assert.NilError(t, err)
	defer conn.Close()

	t.Run("valid", func(t *testing.T) {
		w := db.GormReadWriter{Tx: conn}
		s, err := NewScope(OrganizationScope)
		assert.NilError(t, err)
		assert.Check(t, s.Scope != nil)
		err = w.Create(context.Background(), s)
		assert.NilError(t, err)
		assert.Check(t, s.Id != 0)

		meth, err := NewAuthMethod(s, AuthUserPass)
		assert.NilError(t, err)
		assert.Check(t, meth != nil)
		err = w.Create(context.Background(), meth)
		assert.NilError(t, err)
		assert.Check(t, meth != nil)
		assert.Equal(t, meth.Type, AuthUserPass.String())
	})
	t.Run("nil-scope", func(t *testing.T) {
		meth, err := NewAuthMethod(nil, AuthUserPass)
		assert.Check(t, err != nil)
		assert.Check(t, meth == nil)
		assert.Equal(t, err.Error(), "error scope is nil for new auth method")
	})
}

func TestAuthMethod_GetPrimaryScope(t *testing.T) {
	db.StartTest()
	t.Parallel()
	cleanup, url := db.SetupTest(t, "../db/migrations/postgres")
	defer cleanup()
	defer db.CompleteTest() // must come after the "defer cleanup()"
	conn, err := db.TestConnection(url)
	assert.NilError(t, err)
	defer conn.Close()

	t.Run("valid", func(t *testing.T) {
		w := db.GormReadWriter{Tx: conn}
		s, err := NewScope(OrganizationScope)
		assert.NilError(t, err)
		assert.Check(t, s.Scope != nil)
		err = w.Create(context.Background(), s)
		assert.NilError(t, err)
		assert.Check(t, s.Id != 0)

		meth, err := NewAuthMethod(s, AuthUserPass)
		assert.NilError(t, err)
		assert.Check(t, meth != nil)
		err = w.Create(context.Background(), meth)
		assert.NilError(t, err)
		assert.Check(t, meth != nil)
		assert.Equal(t, meth.Type, AuthUserPass.String())

		primaryScope, err := meth.GetPrimaryScope(context.Background(), &w)
		assert.NilError(t, err)
		assert.Equal(t, primaryScope.GetId(), s.Id)
	})

}

func TestAuthMethod_ResourceType(t *testing.T) {
	db.StartTest()
	t.Parallel()
	cleanup, url := db.SetupTest(t, "../db/migrations/postgres")
	defer cleanup()
	defer db.CompleteTest() // must come after the "defer cleanup()"
	conn, err := db.TestConnection(url)
	assert.NilError(t, err)
	defer conn.Close()

	t.Run("valid", func(t *testing.T) {
		w := db.GormReadWriter{Tx: conn}
		s, err := NewScope(OrganizationScope)
		assert.NilError(t, err)
		assert.Check(t, s.Scope != nil)
		err = w.Create(context.Background(), s)
		assert.NilError(t, err)
		assert.Check(t, s.Id != 0)

		meth, err := NewAuthMethod(s, AuthUserPass)
		assert.NilError(t, err)
		assert.Check(t, meth != nil)
		err = w.Create(context.Background(), meth)
		assert.NilError(t, err)
		assert.Check(t, meth != nil)
		assert.Equal(t, meth.Type, AuthUserPass.String())

		ty := meth.ResourceType()
		assert.Equal(t, ty, ResourceTypeAuthMethod)
	})
}

func TestAuthMethod_Actions(t *testing.T) {
	meth := &AuthMethod{}
	a := meth.Actions()
	assert.Equal(t, a[ActionList.String()], ActionList)
	assert.Equal(t, a[ActionCreate.String()], ActionCreate)
	assert.Equal(t, a[ActionUpdate.String()], ActionUpdate)
	assert.Equal(t, a[ActionEdit.String()], ActionEdit)
	assert.Equal(t, a[ActionDelete.String()], ActionDelete)
}

func TestAuthMethod_Clone(t *testing.T) {
	db.StartTest()
	t.Parallel()
	cleanup, url := db.SetupTest(t, "../db/migrations/postgres")
	defer cleanup()
	defer db.CompleteTest() // must come after the "defer cleanup()"
	conn, err := db.TestConnection(url)
	assert.NilError(t, err)
	defer conn.Close()

	t.Run("valid", func(t *testing.T) {
		w := db.GormReadWriter{Tx: conn}
		s, err := NewScope(OrganizationScope)
		assert.NilError(t, err)
		assert.Check(t, s.Scope != nil)
		err = w.Create(context.Background(), s)
		assert.NilError(t, err)
		assert.Check(t, s.Id != 0)

		meth, err := NewAuthMethod(s, AuthUserPass)
		assert.NilError(t, err)
		assert.Check(t, meth != nil)
		err = w.Create(context.Background(), meth)
		assert.NilError(t, err)
		assert.Check(t, meth != nil)
		assert.Equal(t, meth.Type, AuthUserPass.String())

		cp := meth.Clone()
		assert.Check(t, proto.Equal(cp.(*AuthMethod).AuthMethod, meth.AuthMethod))
	})
	t.Run("not-equal", func(t *testing.T) {
		w := db.GormReadWriter{Tx: conn}
		s, err := NewScope(OrganizationScope)
		assert.NilError(t, err)
		assert.Check(t, s.Scope != nil)
		err = w.Create(context.Background(), s)
		assert.NilError(t, err)
		assert.Check(t, s.Id != 0)

		meth, err := NewAuthMethod(s, AuthUserPass)
		assert.NilError(t, err)
		assert.Check(t, meth != nil)
		err = w.Create(context.Background(), meth)
		assert.NilError(t, err)
		assert.Check(t, meth != nil)
		assert.Equal(t, meth.Type, AuthUserPass.String())

		meth2, err := NewAuthMethod(s, AuthUserPass)
		assert.NilError(t, err)
		assert.Check(t, meth2 != nil)
		err = w.Create(context.Background(), meth2)
		assert.NilError(t, err)

		cp := meth.Clone()
		assert.Check(t, !proto.Equal(cp.(*AuthMethod).AuthMethod, meth2.AuthMethod))
	})
}
