package gen

import (
	"embed"
	"errors"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/cmds/tpm-batis-cli/gen/attribute"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/cmds/tpm-batis-cli/schema"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fileutil"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/templateutil"
	"github.com/rs/zerolog/log"
	godiffpatch "github.com/sourcegraph/go-diff-patch"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/*
var templates embed.FS

const (
	TmplCollectionReadme                 = "templates/%s/readme.txt"
	TmplCollectionEntity                 = "templates/%s/entity.txt"
	TmplCollectionValidateString         = "templates/%s/validate-string.txt"
	TmplCollectionValidateInt            = "templates/%s/validate-int.txt"
	TmplCollectionValidateBool           = "templates/%s/validate-bool.txt"
	TmplCollectionValidateTime           = "templates/%s/validate-time.txt"
	TmplCollectionValidateNullableString = "templates/%s/validate-nullable-string.txt"
	TmplCollectionValidateNullableInt    = "templates/%s/validate-nullable-int.txt"
	TmplCollectionValidateNullableBool   = "templates/%s/validate-nullable-bool.txt"
	TmplCollectionValidateNullableTime   = "templates/%s/validate-nullable-time.txt"
	TmplCollectionEntityTest             = "templates/%s/entity_test.txt"
	TmplCollectionInit                   = "templates/%s/init.txt"
	TmplCollectionMapperXML              = "templates/%s/mapper-xml.txt"
	TmplCollectionUpdate                 = "templates/%s/update.txt"
	TmplCollectionUpdateString           = "templates/%s/update-string.txt"
	TmplCollectionUpdateInt              = "templates/%s/update-int.txt"
	TmplCollectionUpdateBool             = "templates/%s/update-bool.txt"
	TmplCollectionUpdateTime             = "templates/%s/update-time.txt"
	TmplCollectionUpdateNullableString   = "templates/%s/update-nullable-string.txt"
	TmplCollectionUpdateNullableInt      = "templates/%s/update-nullable-int.txt"
	TmplCollectionUpdateNullableBool     = "templates/%s/update-nullable-bool.txt"
	TmplCollectionUpdateNullableTime     = "templates/%s/update-nullable-time.txt"
	TmplCollectionCriteria               = "templates/%s/criteria.txt"
	TmplCollectionCriteriaString         = "templates/%s/criteria-string.txt"
	TmplCollectionCriteriaInt            = "templates/%s/criteria-int.txt"
	TmplCollectionCriteriaBool           = "templates/%s/criteria-bool.txt"
	TmplCollectionCriteriaTime           = "templates/%s/criteria-time.txt"
	TmplCollectionCriteriaNullableString = "templates/%s/criteria-nullable-string.txt"
	TmplCollectionCriteriaNullableInt    = "templates/%s/criteria-nullable-int.txt"
	TmplCollectionCriteriaNullableBool   = "templates/%s/criteria-nullable-bool.txt"
	TmplCollectionCriteriaNullableTime   = "templates/%s/criteria-nullable-time.txt"
	TmplCollectionDelete                 = "templates/%s/delete.txt"
	TmplCollectionInsert                 = "templates/%s/insert.txt"
	TmplCollectionSelect                 = "templates/%s/select.txt"
	TmplCollectionSchema                 = "templates/%s/schema.txt"
)

func entityTmplList(tmplVersion string) []string {
	s := make([]string, 0, 1)
	s = append(s, fmt.Sprintf(TmplCollectionEntity, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionValidateString, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionValidateInt, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionValidateBool, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionValidateTime, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionValidateNullableString, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionValidateNullableInt, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionValidateNullableBool, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionValidateNullableTime, tmplVersion))
	return s
}

func entityTestTmplList(tmplVersion string) []string {
	s := make([]string, 0, 1)
	s = append(s, fmt.Sprintf(TmplCollectionEntityTest, tmplVersion))
	return s
}

func readmeTmplList(tmplVersion string) []string {
	s := make([]string, 0, 1)
	s = append(s, fmt.Sprintf(TmplCollectionReadme, tmplVersion))
	return s
}

func schemaTmplList(tmplVersion string) []string {
	s := make([]string, 0, 1)
	s = append(s, fmt.Sprintf(TmplCollectionSchema, tmplVersion))
	return s
}

func initTmplList(tmplVersion string) []string {
	s := make([]string, 0, 1)
	s = append(s, fmt.Sprintf(TmplCollectionInit, tmplVersion))
	return s
}

func updateTmplList(tmplVersion string) []string {
	s := make([]string, 0, 1)
	s = append(s, fmt.Sprintf(TmplCollectionUpdate, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionUpdateString, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionUpdateInt, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionUpdateBool, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionUpdateTime, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionUpdateNullableString, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionUpdateNullableInt, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionUpdateNullableBool, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionUpdateNullableTime, tmplVersion))
	return s
}

func deleteTmplList(tmplVersion string) []string {
	s := make([]string, 0, 1)
	s = append(s, fmt.Sprintf(TmplCollectionDelete, tmplVersion))
	return s
}

func selectTmplList(tmplVersion string) []string {
	s := make([]string, 0, 1)
	s = append(s, fmt.Sprintf(TmplCollectionSelect, tmplVersion))
	return s
}

func insertTmplList(tmplVersion string) []string {
	s := make([]string, 0, 1)
	s = append(s, fmt.Sprintf(TmplCollectionInsert, tmplVersion))
	return s
}

func criteriaTmplList(tmplVersion string) []string {
	s := make([]string, 0, 1)
	s = append(s, fmt.Sprintf(TmplCollectionCriteria, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionCriteriaString, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionCriteriaInt, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionCriteriaBool, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionCriteriaTime, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionCriteriaNullableString, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionCriteriaNullableInt, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionCriteriaNullableBool, tmplVersion))
	s = append(s, fmt.Sprintf(TmplCollectionCriteriaNullableTime, tmplVersion))
	return s
}

func mapperXMLTmpList(tmplVersion string) []string {
	s := make([]string, 0, 1)
	s = append(s, fmt.Sprintf(TmplCollectionMapperXML, tmplVersion))
	return s
}

type Options struct {
	TargetFolder string
	Version      string
	FormatCode   bool
}

type Option func(*Options)

func WithOutputFolder(f string) Option {
	return func(opts *Options) {
		opts.TargetFolder = f
	}
}

func WithVersion(v string) Option {
	return func(opts *Options) {
		opts.Version = v
	}
}

func WithFormatCode(b bool) Option {
	return func(opts *Options) {
		opts.FormatCode = b
	}
}

// GenerationContext is the object passed to templates
type GenerationContext struct {
	Opts         Options
	Schema       *schema.Schema
	Attributes   []attribute.Attribute
	PkAttributes []attribute.Attribute
}

func (ctx GenerationContext) GoPackageImports(ambit string) []string {
	pkgs := make([]string, 0)
	for _, a := range ctx.Attributes {
		switch ambit {
		case "entity":
			if ap := a.GoPackageImports(); len(ap) > 0 {
				pkgs = append(pkgs, ap...)
			}
		case "criteria":
			if a.GetDefinition().WithCriterion {
				if ap := a.GoPackageImports(); len(ap) > 0 {
					pkgs = append(pkgs, ap...)
				}
			}
		case "update":
			// Should put some flag..... not all fields go to the update stat
			if !a.GetDefinition().IsPKey {
				if ap := a.GoPackageImports(); len(ap) > 0 {
					pkgs = append(pkgs, ap...)
				}
			}
		}
	}

	if len(pkgs) > 1 {
		pkgs = util.RemoveStringDuplicates(pkgs)
	}

	return pkgs
}

func (ctx GenerationContext) MaxTextType(nullable bool) []int {
	maxTextTypesMap := make(map[int]struct{})
	var res []int

	for _, a := range ctx.Attributes {
		if (a.GetDefinition().Typ == "string" && a.GetDefinition().Nullable == false) ||
			(a.GetDefinition().Typ == "nullable-string" && a.GetDefinition().Nullable == true) {

			if _, ok := maxTextTypesMap[a.GetDefinition().MaxLength]; !ok {
				res = append(res, a.GetDefinition().MaxLength)
				maxTextTypesMap[a.GetDefinition().MaxLength] = struct{}{}
			}
		}
	}

	/*
		if len(maxTextTypesMap) > 0 {
			var res []int
			for i, _ := range maxTextTypesMap {
				res = append(res, i)
			}

			return res
		}
	*/

	return res
}

type Generator struct {
	Opts       Options
	Schema     *schema.Schema
	Attributes []attribute.Attribute
}

func NewGenerator(sch *schema.Schema, opts ...Option) (Generator, error) {
	const semLogContext = "sql-cli-generator::new"

	g := Generator{Schema: sch}

	g.Opts = Options{Version: "v1"}
	for _, o := range opts {
		o(&g.Opts)
	}

	for _, a := range sch.Fields {
		ga, err := attribute.NewAttribute(a)
		if err != nil {
			return g, err
		}
		g.Attributes = append(g.Attributes, ga)
	}

	return g, nil
}

func (g *Generator) PkAttributes() []attribute.Attribute {
	var pk []attribute.Attribute
	for _, a := range g.Attributes {
		if a.GetDefinition().IsPKey {
			pk = append(pk, a)
		}
	}

	return pk
}

func (g *Generator) Generate() error {
	const semLogContext = "sql-cli-generator::generate"
	log.Info().Str("target-folder", g.Opts.TargetFolder).Msg(semLogContext + " starting")

	genCtx := GenerationContext{g.Opts, g.Schema, g.Attributes, g.PkAttributes()}

	if err := g.emit(genCtx, g.Opts.TargetFolder, "readme.md", readmeTmplList(g.Opts.Version), false); err != nil {
		return err
	}

	if err := g.emit(genCtx, g.Opts.TargetFolder, "schema.go", schemaTmplList(g.Opts.Version), false); err != nil {
		return err
	}

	if err := g.emit(genCtx, g.Opts.TargetFolder, "init.go", initTmplList(g.Opts.Version), g.Opts.FormatCode); err != nil {
		return err
	}

	if err := g.emit(genCtx, g.Opts.TargetFolder, "mapper.xml", mapperXMLTmpList(g.Opts.Version), false); err != nil {
		return err
	}

	if err := g.emit(genCtx, g.Opts.TargetFolder, "entity.go", entityTmplList(g.Opts.Version), g.Opts.FormatCode); err != nil {
		return err
	}

	if g.Schema.Properties.DbType == "table" {
		if err := g.emit(genCtx, g.Opts.TargetFolder, "update.go", updateTmplList(g.Opts.Version), g.Opts.FormatCode); err != nil {
			return err
		}
	}

	if err := g.emit(genCtx, g.Opts.TargetFolder, "criteria.go", criteriaTmplList(g.Opts.Version), g.Opts.FormatCode); err != nil {
		return err
	}

	if g.Schema.Properties.DbType == "table" {
		if err := g.emit(genCtx, g.Opts.TargetFolder, "delete.go", deleteTmplList(g.Opts.Version), g.Opts.FormatCode); err != nil {
			return err
		}
	}

	if g.Schema.Properties.DbType == "table" {
		if err := g.emit(genCtx, g.Opts.TargetFolder, "insert.go", insertTmplList(g.Opts.Version), g.Opts.FormatCode); err != nil {
			return err
		}
	}

	if err := g.emit(genCtx, g.Opts.TargetFolder, "select.go", selectTmplList(g.Opts.Version), g.Opts.FormatCode); err != nil {
		return err
	}

	if err := g.emit(genCtx, g.Opts.TargetFolder, "entity_test.go", entityTestTmplList(g.Opts.Version), g.Opts.FormatCode); err != nil {
		return err
	}

	return nil
}

func (g *Generator) emit(ctx GenerationContext, outFolder string, generatedFileName string, templates []string, formatCode bool) error {
	if t, ok := loadTemplate(templates...); ok {
		destinationFile := filepath.Join(outFolder, generatedFileName)
		log.Info().Str("dest", destinationFile).Msgf("generating text from template")

		if err := parseTemplateWithFuncMapsProcessWrite2File(t, getTemplateUtilityFunctions(), ctx, destinationFile, formatCode); err != nil {
			log.Error().Err(err).Msgf("parse template failed")
			return err
		}
	} else {
		log.Info().Msgf("unable to load template ...skipping")
		return errors.New("unable to load template ...skipping")
	}

	return nil
}

func loadTemplate(templatePath ...string) ([]templateutil.Info, bool) {

	res := make([]templateutil.Info, 0)
	for _, tpath := range templatePath {

		/*
		 * Get the name of the template from the file name.... Hope there is one dot only...
		 * Dunno it is a problem.
		 */
		tname := filepath.Base(tpath)
		if ext := filepath.Ext(tname); ext != "" {
			tname = strings.TrimSuffix(tname, ext)
		}

		var b []byte
		var e error

		b, e = templates.ReadFile(tpath)
		if e != nil {
			return nil, false
		}

		res = append(res, templateutil.Info{Name: tname, Content: string(b)})
	}

	return res, true
}

func parseTemplateWithFuncMapsProcessWrite2File(templates []templateutil.Info, fMaps template.FuncMap, templateData interface{}, outputFile string, formatSource bool) error {

	if pkgTemplate, err := templateutil.Parse(templates, fMaps); err != nil {
		return err
	} else {
		currentContent, ok, err := readCurrentFileContent(outputFile)
		if err != nil {
			return err
		}

		if err := templateutil.ProcessWrite2File(pkgTemplate, templateData, outputFile, formatSource); err != nil {
			return err
		}

		newContent, _, err := readCurrentFileContent(outputFile)
		if err != nil {
			return err
		}

		if ok {
			patch := godiffpatch.GeneratePatch(outputFile, string(currentContent), string(newContent))
			if len(patch) > 0 {
				patchFile := filepath.Join(filepath.Dir(outputFile), filepath.Base(outputFile)+".patch")
				_ = os.WriteFile(patchFile, []byte(patch), fs.ModePerm)
			}
		}
	}

	return nil
}

func readCurrentFileContent(fn string) ([]byte, bool, error) {
	if fileutil.FileExists(fn) {
		b, err := os.ReadFile(fn)
		if err != nil {
			return nil, true, err
		}

		return b, true, nil
	}

	return nil, false, nil
}
func getTemplateUtilityFunctions() template.FuncMap {

	fMap := template.FuncMap{
		"format": func(n string, mode string) string {
			r := n
			switch mode {
			case "capitalize":
				r = util.Capitalize(n)
			case "to-upper":
				r = strings.ToUpper(n)
			case "to-lower":
				r = strings.ToLower(n)
			}

			return r
		},
		"filterSubTemplateContext": func(attribute attribute.Attribute, criteriaObjectName string) map[string]interface{} {
			return map[string]interface{}{
				"Attr":              attribute,
				"CriteriaStructRef": criteriaObjectName,
			}
		},
	}

	return fMap
}
