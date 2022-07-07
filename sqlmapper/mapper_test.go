package sqlmapper

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/system/util"

	"testing"
)

var whereIfXML01 = []byte(`
<sqlmapper namespace="org.r3.db.system.site.SiteMapper" >

    <resultMap id="BaseResultMap" type="org.r3.db.system.site.SiteDTO" >
         <id column="site" property="site" jdbcType="VARCHAR" />
           <result column="parentsite" property="parentsite" jdbcType="VARCHAR" />
           <result column="sitedescr" property="sitedescr" jdbcType="VARCHAR" />
           <collection 	property="list_properties" 
          				column="propertykey1"  
          				>           
          		<id column="propertykey1" property="propertykey1" jdbcType="VARCHAR" />
          		<result column="propertyscope1" property="propertyscope1" jdbcType="VARCHAR" />        
                <collection property="list_list_properties" 
          				column="propertykey2" 
          				javaType="org.r3.db.system.siteproperty.SitePropertyDTOCollection" 
          				>           
          		   <id column="propertykey2" property="propertykey2" jdbcType="VARCHAR" />
          		   <result column="propertyscope2" property="propertyscope2" jdbcType="VARCHAR" />           
                </collection>

          </collection>
    </resultMap>

  <select id="selectBySite" resultMap="BaseResultMap"  parameterType="map" >
    select * from ${.tableName}
    <where> 
      <trim  prefixOverrides="and" >
      <if test=".field1" >
         and f1 = #{.field1,jdbcType=VARCHAR}
      </if>

      <if test=".field2" >
         and f2 = #{.field2,jdbcType=VARCHAR}
      </if>

      <if test=".field3" >
         and f3 = #{.field3,jdbcType=VARCHAR}
      </if>
      </trim>
 
    </where>
  </select>

</sqlmapper>
`)

func TestWhereIf01(t *testing.T) {

	var wanted string
	var sqlStmt MappedStatement

	if m, err := NewMapper(string(whereIfXML01), WithBindStyle(BINDSTYLE_QUESTION)); err != nil {
		t.Fatal(err)
	} else {

		wanted = "select * from r3_site where f1 = ? and f2 = ?"
		mapp := map[string]interface{}{
			"field1":    "field1_value",
			"field2":    "field2_value",
			"tableName": "r3_site",
		}
		sqlStmt, err = m.GetMappedStatement("selectBySite", mapp)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("Statement", sqlStmt.GetStatement(), "Params", fmt.Sprint(sqlStmt.GetParams()))
		if sqlStmt.GetStatement() != wanted {
			t.Error(fmt.Sprintf("Result doesn't match wanted: result = %q wanted: %q", sqlStmt.GetStatement(), wanted))
		}

		wanted = "select * from r3_site"
		mapp = map[string]interface{}{
			"tableName": "r3_site",
		}
		sqlStmt, err = m.GetMappedStatement("selectBySite", mapp)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("Statement", sqlStmt.GetStatement(), "Params", fmt.Sprint(sqlStmt.GetParams()))
		if sqlStmt.GetStatement() != wanted {
			t.Error(fmt.Sprintf("Result doesn't match wanted: result = %q wanted: %q", sqlStmt.GetStatement(), wanted))
		}

		wanted = "select * from r3_site where f2 = ?"
		mapp = map[string]interface{}{
			"tableName": "r3_site",
			"field2":    "field2_value",
		}
		sqlStmt, err = m.GetMappedStatement("selectBySite", mapp)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("Statement", sqlStmt.GetStatement(), "Params", fmt.Sprint(sqlStmt.GetParams()))
		if sqlStmt.GetStatement() != wanted {
			t.Error(fmt.Sprintf("Result doesn't match wanted: result = %q wanted: %q", sqlStmt.GetStatement(), wanted))
		}

	}
}

var whereChooseXML01 = []byte(`
<sqlmapper namespace="org.r3.db.system.site.SiteMapper" >

  <select id="selectBySite" resultMap="BaseResultMap" parameterType="map" >
    select * from ${.tableName}
    <where> 

      <trim prefix="(" suffix=")" prefixOverrides="or" >
	  <choose>
      <when test=".field4" >
         or f4 = #{.field4,jdbcType=VARCHAR}
      </when>

      <when test=".field5" >
         or f5 = #{.field5,jdbcType=VARCHAR}
      </when>

      <otherwise>
         or otherwise...
      </otherwise>
      </choose>      
      </trim>

    </where>
  </select>

</sqlmapper>
`)

func TestWhereChoose01(t *testing.T) {

	var sqlStmt MappedStatement
	var wanted string

	if m, err := NewMapper(string(whereChooseXML01), WithBindStyle(BINDSTYLE_QUESTION)); err != nil {
		t.Fatal(err)
	} else {

		wanted = "select * from r3_site where (f5 = ?)"
		mapp := map[string]interface{}{
			"field5":    "field5_value",
			"tableName": "r3_site",
		}
		sqlStmt, err = m.GetMappedStatement("selectBySite", mapp)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("Statement", sqlStmt.GetStatement(), "Params", fmt.Sprint(sqlStmt.GetParams()))
		if sqlStmt.GetStatement() != wanted {
			t.Error(fmt.Sprintf("Result doesn't match wanted: result = %q wanted: %q", sqlStmt.GetStatement(), wanted))
		}

		wanted = "select * from r3_site where (otherwise...)"
		mapp = map[string]interface{}{
			"tableName": "r3_site",
		}
		sqlStmt, err = m.GetMappedStatement("selectBySite", mapp)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("Statement", sqlStmt.GetStatement(), "Params", fmt.Sprint(sqlStmt.GetParams()))
		if sqlStmt.GetStatement() != wanted {
			t.Error(fmt.Sprintf("Result doesn't match wanted: result = %q wanted: %q", sqlStmt.GetStatement(), wanted))
		}

	}
}

var paramPlaceHolderXML = []byte(`
<sqlmapper namespace="org.r3.db.system.site.SiteMapper" >

  <select id="selectBySite" resultMap="BaseResultMap" parameterType="map" >
    select * from ${.tableName}
    <where> 
      <trim  prefixOverrides="and" >
      <if test=".field1" >
         and f1 = #{.field1,jdbcType=VARCHAR}
      </if>
      </trim>
    </where>
  </select>

</sqlmapper>
`)

func TestParamStyle(t *testing.T) {

	var sqlStmt MappedStatement
	var wanted string

	wanted = "select * from r3_site where f1 = ?"
	if m, err := NewMapper(string(paramPlaceHolderXML), WithBindStyle(BINDSTYLE_QUESTION)); err != nil {
		t.Fatal(err)
	} else {
		mapp := map[string]interface{}{
			"field1":    "field1_value",
			"tableName": "r3_site",
		}
		sqlStmt, err = m.GetMappedStatement("selectBySite", mapp)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("Statement", sqlStmt.GetStatement(), "Params", fmt.Sprint(sqlStmt.GetParams()))
		if sqlStmt.GetStatement() != wanted {
			t.Error(fmt.Sprintf("Result doesn't match wanted: result = %q wanted: %q", sqlStmt.GetStatement(), wanted))
		}
	}

	wanted = "select * from r3_site where f1 = S1"
	if m, err := NewMapper(string(paramPlaceHolderXML), WithBindStyle(BINDSTYLE_DOLLAR)); err != nil {
		t.Fatal(err)
	} else {
		mapp := map[string]interface{}{
			"field1":    "field1_value",
			"tableName": "r3_site",
		}
		sqlStmt, err = m.GetMappedStatement("selectBySite", mapp)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("Statement", sqlStmt.GetStatement(), "Params", fmt.Sprint(sqlStmt.GetParams()))
		if sqlStmt.GetStatement() != wanted {
			t.Error(fmt.Sprintf("Result doesn't match wanted: result = %q wanted: %q", sqlStmt.GetStatement(), wanted))
		}
	}

}

var filterCriteriaXML = []byte(`
<sqlmapper namespace='org.r3.db.system.site.SiteMapper' >

  <sql id='Filter_Where_Clause' >
    <where >
      <foreach collection='.filter.OrListOfCriteria' item='criteria' separator='or' >
          <trim prefix='(' suffix=')' prefixOverrides='and' >
            <foreach collection='.criteria.AndListOfCriterion' item='criterion' >
              <choose>
                <when test='eq .criterion.Type "NoValue"' >
                  and ${.criterion.Condition}
                </when>
                <when test='eq .criterion.Type "SingleValue"' >
                  and ${.criterion.Condition} #{.criterion.Value}
                </when>
                <when test='eq .criterion.Type "BetweenValue"' >
                  and ${.criterion.Condition} #{.criterion.Value} and #{.criterion.secondValue}
                </when>
                <when test='eq .criterion.Type "ListValue"' >
                  and ${.criterion.Condition}
                  <foreach collection='.criterion.Value' item='listItem' open='(' close=')' separator=',' >
                    #{.listItem}
                  </foreach>
                </when> 
              </choose>
            </foreach>
          </trim>
      </foreach>
    </where>
  </sql>

  <select id='selectBySite' resultMap='BaseResultMap' parameterType='map' >
    select * from ${.tableName}
    <include refid='Filter_Where_Clause' />    
  </select>

</sqlmapper>
`)

type TestFilterBuilder struct {
	fb *FilterBuilder
}

func NewTestFilterBuilder() *TestFilterBuilder {
	return &TestFilterBuilder{fb: &FilterBuilder{}}
}

func (ub *TestFilterBuilder) Or() *TestFilterBuilder {
	ub.fb.Or()
	return ub
}

func (ub *TestFilterBuilder) Build() Filter {
	return ub.fb.Build()
}

func (ub *TestFilterBuilder) AndCustomCondition(condition string) *TestFilterBuilder {
	ub.fb.And(Criterion{Type: NoValue, Condition: condition})
	return ub
}

func (ub *TestFilterBuilder) AndSiteEqualTo(site string) *TestFilterBuilder {
	ub.fb.And(Criterion{Type: SingleValue, Condition: "site = ", Value: site})
	return ub
}

func (ub *TestFilterBuilder) AndParentSiteIn(sites []string) *TestFilterBuilder {
	ub.fb.And(Criterion{Type: ListValue, Condition: "parentSite in ", Value: sites})
	return ub
}

func TestFilterCriteria(t *testing.T) {

	var sqlStmt MappedStatement
	var wanted string

	wanted = "select * from r3_site where f1 = ?"
	if m, err := NewMapper(string(filterCriteriaXML), WithBindStyle(BINDSTYLE_QUESTION)); err != nil {
		t.Fatal(err)
	} else {

		f := NewFilterBuilder().Or()
		f.And(Criterion{Type: NoValue, Condition: "My_Condition = 71"})
		f.And(Criterion{Type: SingleValue, Condition: "site = ", Value: "site"})
		f.And(Criterion{Type: SingleValue, Condition: "parentSite in", Value: [3]string{"s1", "s2", "s3"}})

		var mapp map[string]interface{}
		mapp = map[string]interface{}{
			"field1":    "field1_value",
			"tableName": "r3_site",
			"filter":    f.Build(),
		}

		wanted = "select * from r3_site where (My_Condition = 71 and site = ? and parentSite in (?, ?, ?))"
		sqlStmt, err = m.GetMappedStatement("selectBySite", mapp)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("Statement", sqlStmt.GetStatement(), "Params", fmt.Sprint(sqlStmt.GetParams()))
		if sqlStmt.GetStatement() != wanted {
			t.Error(fmt.Sprintf("Result doesn't match wanted: result = %q wanted: %q", sqlStmt.GetStatement(), wanted))
		}

		f = NewFilterBuilder()
		f.Or().And(Criterion{Type: NoValue, Condition: "My_Condition = 71"})
		f.Or().And(Criterion{Type: SingleValue, Condition: "site = ", Value: "site"})
		f.Or().And(Criterion{Type: SingleValue, Condition: "parentSite in", Value: [3]string{"s1", "s2", "s3"}})

		mapp = map[string]interface{}{
			"tableName": "r3_site",
			"filter":    f.Build(),
		}
		wanted = "select * from r3_site where (My_Condition = 71) or (site = ?) or (parentSite in (?, ?, ?))"
		sqlStmt, err = m.GetMappedStatement("selectBySite", mapp)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("Statement", sqlStmt.GetStatement(), "Params", fmt.Sprint(sqlStmt.GetParams()))
		if sqlStmt.GetStatement() != wanted {
			t.Error(fmt.Sprintf("Result doesn't match wanted: result = %q wanted: %q", sqlStmt.GetStatement(), wanted))
		}

		/*
		 * Costruzione di un Filter ipotizzando delle API piu' immediate.
		 */
		f1 := NewTestFilterBuilder()
		f1.Or().AndCustomCondition("field is not null")
		f1.Or().AndSiteEqualTo("site")
		f1.Or().AndParentSiteIn([]string{"s4", "s5", "s6"})

		mapp = map[string]interface{}{
			"tableName": "r3_site",
			"filter":    f1.Build(),
		}
		wanted = "select * from r3_site where (field is not null) or (site = ?) or (parentSite in ( ? , ? , ? ))"
		sqlStmt, err = m.GetMappedStatement("selectBySite", mapp)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("Statement", sqlStmt.GetStatement(), "Params", fmt.Sprint(sqlStmt.GetParams()))
		if sqlStmt.GetStatement() != wanted {
			t.Error(fmt.Sprintf("Result doesn't match wanted: result = %q wanted: %q", sqlStmt.GetStatement(), wanted))
		}

	}

}

func TestFilterAPI(t *testing.T) {

	// f := Filter{}
	f := NewFilterBuilder().Or().Or().And(Criterion{Type: NoValue, Condition: "My_Condition = 71"}).And(Criterion{Type: SingleValue, Condition: "site = ", Value: "site"})
	fmt.Println(f.Build())

	f = NewFilterBuilder().And(Criterion{Type: NoValue, Condition: "My_Condition = 71"}).And(Criterion{Type: SingleValue, Condition: "site = ", Value: "site"})
	fmt.Println(f.Build())

	f = NewFilterBuilder()
	fmt.Println(f.Build())
}

var configMapper = []byte(`
<sqlmapper namespace="org.r3.db.test_mapper" >

  <select id="selectBySite" resultMap="BaseResultMap" >
    select s.site from r3_site s
    where site = #{.site,jdbcType=VARCHAR}
  </select>

  <select id="selectByUser" resultMap="BaseResultMap" >
    select * from r3_user s
    where nickname = #{.nickname,jdbcType=VARCHAR}
  </select>

</sqlmapper>
`)

var configFileResource = []byte(`
<configuration>
  <mappers>
     <sqlmapper resource="resources/sqlmapper.xml"/>
  </mappers>
</configuration> 
`)

func TestMapperRegistry(t *testing.T) {

	filn, ln, fun := util.GetExecutingFunctionInfo()
	t.Log(fmt.Sprintf("%s:%d %s\n", filn, ln, fun))

	registry := make(map[string][]byte)
	registry["resources/goBatisCfg.xml"] = configFileResource
	registry["resources/sqlmapper.xml"] = configMapper
	_, err := NewRegistry(util.NewFileRegistryCascadeResolver(registry), "resources/goBatisCfg.xml")
	if err != nil {
		t.Fatal(err)
	}

}
