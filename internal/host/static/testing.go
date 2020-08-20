package static

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/boundary/internal/db"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

// TestCatalogs creates count number of static host catalogs to the provided DB
// with the provided scope id.  If any errors are encountered during the creation of
// the host catalog, the test will fail.
func TestCatalogs(t *testing.T, conn *gorm.DB, scopeId string, count int) []*HostCatalog {
	t.Helper()
	assert := assert.New(t)
	var cats []*HostCatalog
	for i := 0; i < count; i++ {
		cat, err := NewHostCatalog(scopeId)
		assert.NoError(err)
		assert.NotNil(cat)
		id, err := newHostCatalogId()
		assert.NoError(err)
		assert.NotEmpty(id)
		cat.PublicId = id

		w := db.New(conn)
		err2 := w.Create(context.Background(), cat)
		assert.NoError(err2)
		cats = append(cats, cat)
	}
	return cats
}

// TestHosts creates count number of static hosts to the provided DB
// with the provided catalog id.  The catalog must have been created previously.
// If any errors are encountered during the creation of the host, the test will fail.
func TestHosts(t *testing.T, conn *gorm.DB, catalogId string, count int) []*Host {
	t.Helper()
	assert := assert.New(t)
	var hosts []*Host

	for i := 0; i < count; i++ {
		host, err := NewHost(catalogId, WithAddress(fmt.Sprintf("%s-%d", catalogId, i)))
		assert.NoError(err)
		assert.NotNil(host)

		id, err := newHostId()
		assert.NoError(err)
		assert.NotEmpty(id)
		host.PublicId = id

		w := db.New(conn)
		err2 := w.Create(context.Background(), host)
		assert.NoError(err2)
		hosts = append(hosts, host)
	}
	return hosts
}