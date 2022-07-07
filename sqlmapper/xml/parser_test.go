package xml

import (
	"testing"
)

var data = []byte(`
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

  <sql id="Filter_Where_Clause" >
    <where >
      <foreach collection="filter.oredCriteria" item="criteria" separator="or" >
        <if test="criteria.valid" >
          <trim prefix="(" suffix=")" prefixOverrides="and" >
            <foreach collection="criteria.criteria" item="criterion" >
              <choose >
                <when test="criterion.noValue" >
                  and ${criterion.condition}
                </when>
                <when test="criterion.singleValue" >
                  and ${criterion.condition} #{criterion.value}
                </when>
                <when test="criterion.betweenValue" >
                  and ${criterion.condition} #{criterion.value} and #{criterion.secondValue}
                </when>
                <when test="criterion.listValue" >
                  and ${criterion.condition}
                  <foreach collection="criterion.value" item="listItem" open="(" close=")" separator="," >
                    #{listItem}
                  </foreach>
                </when>
                <when test="criterion.subSelectCriterion" >
                  and ${criterion.prefixCondition} #{criterion.value} ${criterion.suffixCondition}
                </when>                                
              </choose>
            </foreach>
          </trim>
        </if>
      </foreach>
    </where>
  </sql>

  <sql id="Base_Column_List" >
    site,  parentsite,  sitedescr,  sitelanguages,  sitesummary,  sitedomain,  calendarid,  hostname,  homeurl,
  consoleurl,  sitetype  </sql>
  
  <sql id='GroupBy_Column_List'>
    <trim suffixOverrides=",">
    <choose>
       <when test="filter != null &amp;&amp; filter.isGroupBy('site')" >
           site,
       </when>
       <otherwise>
           null as site,
       </otherwise>         
     </choose> 
     </trim>       
  </sql>
  
  <select id="selectBySite" resultMap="BaseResultMap" parameterType="map" >
    select 
    <include refid="Base_Column_List" />
    from r3_site
    where site = #{site,jdbcType=VARCHAR}
  </select>
  <delete id="deleteBySite" parameterType="map" >
    delete
    from r3_site
    where site = #{site,jdbcType=VARCHAR}
  </delete>
  <select id="selectByPrimaryKey" resultMap="BaseResultMap" parameterType="map" >
    select 
    <include refid="Base_Column_List" />
    from r3_site    
    where  site = #{pk.site,jdbcType=VARCHAR}
  </select>
  <delete id="deleteByPrimaryKey" parameterType="map">
    delete
    from r3_site
    where  site = #{pk.site,jdbcType=VARCHAR}
  </delete>  
  <select id="select" resultMap="BaseResultMap" parameterType="map" >
    select
    <if test="filter != null &amp;&amp; filter.distinct" >
      distinct
    </if>
    
    <choose>
       <when test="filter != null &amp;&amp; filter.groupByClause != null" >
           <include refid="GroupBy_Column_List" />
       </when>
       <otherwise>
           <include refid="Base_Column_List" />
       </otherwise>
    </choose>
        
    from r3_site
    <if test="filter != null" >
      <include refid="Filter_Where_Clause" />
    </if>
    <if test="filter != null &amp;&amp; filter.groupByClause != null" >
      group by ${filter.groupByClause}
    </if>
    <if test="filter != null &amp;&amp; filter.orderByClause != null" >
      order by ${filter.orderByClause}
    </if>
    <if test="filter != null &amp;&amp; filter.limitFetch > 0" >
      limit ${filter.limitFetch}
    </if>
  </select>
  <select id="count" parameterType="map" resultType="java.lang.Integer" >
    select count(*) from r3_site
    <if test="filter != null" >
      <include refid="Filter_Where_Clause" />
    </if>
  </select>  
  <delete id="delete" parameterType="map" >
    delete from r3_site
    <!-- _parameter: dato che e' una map lo sostituisco con il nome del parametro in quanto non è l'unico -->
    <if test="filter != null" >
      <include refid="Filter_Where_Clause" />
    </if>
  </delete>
  <insert id="insert" parameterType="map" >
    insert into r3_site (
 site, parentsite, sitedescr, sitelanguages, sitesummary, sitedomain, calendarid, hostname, homeurl,
 consoleurl, sitetype      )
    values (
 
 #{record.site,jdbcType=VARCHAR}, 
 #{record.parentsite,jdbcType=VARCHAR}, 
 #{record.sitedescr,jdbcType=VARCHAR}, 
 #{record.sitelanguages,jdbcType=VARCHAR},
 
 #{record.sitesummary,jdbcType=VARCHAR}, 
 #{record.sitedomain,jdbcType=VARCHAR}, 
 #{record.calendarid,jdbcType=VARCHAR}, 
 #{record.hostname,jdbcType=VARCHAR},
 
 #{record.homeurl,jdbcType=VARCHAR}, 
 #{record.consoleurl,jdbcType=VARCHAR}, 
 #{record.sitetype,jdbcType=VARCHAR}
      )
  </insert>
  <update id="update" parameterType="map" >
    update r3_site
    <set >
      <if test="record.parentsiteDirty" >
        parentsite = #{record.parentsite,jdbcType=VARCHAR},
      </if>      
      <if test="record.sitedescrDirty" >
        sitedescr = #{record.sitedescr,jdbcType=VARCHAR},
      </if>      
      <if test="record.sitelanguagesDirty" >
        sitelanguages = #{record.sitelanguages,jdbcType=VARCHAR},
      </if>      
      <if test="record.sitesummaryDirty" >
        sitesummary = #{record.sitesummary,jdbcType=VARCHAR},
      </if>      
      <if test="record.sitedomainDirty" >
        sitedomain = #{record.sitedomain,jdbcType=VARCHAR},
      </if>      
      <if test="record.calendaridDirty" >
        calendarid = #{record.calendarid,jdbcType=VARCHAR},
      </if>      
      <if test="record.hostnameDirty" >
        hostname = #{record.hostname,jdbcType=VARCHAR},
      </if>      
      <if test="record.homeurlDirty" >
        homeurl = #{record.homeurl,jdbcType=VARCHAR},
      </if>      
      <if test="record.consoleurlDirty" >
        consoleurl = #{record.consoleurl,jdbcType=VARCHAR},
      </if>      
      <if test="record.sitetypeDirty" >
        sitetype = #{record.sitetype,jdbcType=VARCHAR},
      </if>      
    </set>
    <if test="filter != null" >
      <include refid="Filter_Where_Clause" />
      <!-- Non dovrebbe servire in quanto il filter e' sempre named. -->
      <!-- <include refid="Filter_Where_Clause_4_Update" /> -->
    </if>
  </update>
  <update id="updateByPrimaryKey" parameterType="map" >
    update r3_site
    <set >
      <if test="record.parentsiteDirty" >
        parentsite = #{record.parentsite,jdbcType=VARCHAR},
      </if>
      <if test="record.sitedescrDirty" >
        sitedescr = #{record.sitedescr,jdbcType=VARCHAR},
      </if>
      <if test="record.sitelanguagesDirty" >
        sitelanguages = #{record.sitelanguages,jdbcType=VARCHAR},
      </if>
      <if test="record.sitesummaryDirty" >
        sitesummary = #{record.sitesummary,jdbcType=VARCHAR},
      </if>
      <if test="record.sitedomainDirty" >
        sitedomain = #{record.sitedomain,jdbcType=VARCHAR},
      </if>
      <if test="record.calendaridDirty" >
        calendarid = #{record.calendarid,jdbcType=VARCHAR},
      </if>
      <if test="record.hostnameDirty" >
        hostname = #{record.hostname,jdbcType=VARCHAR},
      </if>
      <if test="record.homeurlDirty" >
        homeurl = #{record.homeurl,jdbcType=VARCHAR},
      </if>
      <if test="record.consoleurlDirty" >
        consoleurl = #{record.consoleurl,jdbcType=VARCHAR},
      </if>
      <if test="record.sitetypeDirty" >
        sitetype = #{record.sitetype,jdbcType=VARCHAR},
      </if>
    </set>
    where  site = #{record.site,jdbcType=VARCHAR}
  </update>

</sqlmapper>
`)

func TestParser_001(t *testing.T) {

	melement, err := ParseXML(string(data))
	if err != nil {
		t.Fatal(err)
	}

	melement.Walk(0, WithTreePrintWalkOption())
}

var data2 = []byte(`
<sqlmapper namespace="org.r3.db.system.site.SiteMapper" >

  <select id="selectBySite" resultMap="BaseResultMap" parameterType="map" >
    select 
    <include refid="Base_Column_List" >
        <property name="pbase" value="pbasevalue" /> 
    </include>
    from r3_site
    where site = #{site,jdbcType=VARCHAR}
  </select>

<sql id="Base_Column_List" >
    site,  ${pbase}, parentsite,  sitedescr,  sitelanguages,  sitesummary,  sitedomain,  calendarid,  hostname,  homeurl,  consoleurl,  sitetype  
    <include refid="Nested_Base_Column_List" >
        <property name="pnested" value="pnestedvalue" /> 
    </include>
  </sql>

  <sql id="Nested_Base_Column_List" >
    Nested, sql, text ${pnested} </sql>

</sqlmapper>
`)

func TestParser_002(t *testing.T) {

	mapper, err := ParseMapper(string(data2))
	if err != nil {
		t.Fatal(err)
	}

	mapper.Walk(0, WithTreePrintWalkOption())
}

var data3 = []byte(`
<sqlmapper namespace="org.r3.db.system.site.SiteMapper" >

  <select id="selectBySite" resultMap="BaseResultMap" parameterType="map" >
    select * from r3_site
    where site = #{site,jdbcType=VARCHAR}
  </select>

</sqlmapper>
`)

func TestParser_003(t *testing.T) {

	mapper, err := ParseMapper(string(data3))
	if err != nil {
		t.Fatal(err)
	}

	mapper.Walk(0, WithTreePrintWalkOption())
}

func TestStatementParam_001(t *testing.T) {
	vmatcher, err := NewStatementParamParser()
	if err != nil {
		panic(err)
	}

	s := "owqcsn ${filter.condition,_,mustexists} poakkdsnvkdsn ${filter.condition3} #{.filter.criteria}"
	vars, _ := vmatcher.GetDollarVariables(s)
	t.Log("Dollar: ", vars)

	vars, _ = vmatcher.GetDashVariables(s)
	t.Log("Dash: ", vars)
}
