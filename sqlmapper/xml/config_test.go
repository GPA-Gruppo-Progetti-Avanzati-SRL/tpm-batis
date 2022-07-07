package xml

import (
	"errors"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/system/util"
	"testing"
)

var mapperXML01 = []byte(`
<sqlmapper namespace="org.r3.db.system.site.SiteMapper" >

  <select id="selectBySite" resultMap="BaseResultMap" parameterType="map" >
    select * from r3_site
    where site = #{site,jdbcType=VARCHAR}
  </select>

</sqlmapper>
`)

var cfgXML01 = []byte(`
<configuration>
  <mappers>
     <sqlmapper resource="sqlmapper.xml"/>
  </mappers>
</configuration> 
`)

func TestConfig(t *testing.T) {

	resolver := func(resourceName string) ([]byte, error) {
		if resourceName == "gobatiscfg.xml" {
			return cfgXML01, nil
		}

		if resourceName == "sqlmapper.xml" {
			return mapperXML01, nil
		}

		return nil, errors.New(fmt.Sprintf("ResourceNotFound: %s", resourceName))
	}

	cfg, err := NewConfig(util.ResourceResolverFunc(resolver), "gobatiscfg.xml")
	if err != nil {
		t.Fatal(err)
	}

	if len(cfg.ListOfMappers.Mappers) != 1 {
		t.Fatal(err)
	}

}
