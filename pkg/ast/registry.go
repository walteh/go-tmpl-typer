package ast

import (
	"context"
	"go/types"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rs/zerolog"
	"gitlab.com/tozd/go/errors"
	"golang.org/x/tools/go/packages"
)

// Registry manages Go package and type information
type Registry struct {
	// Types maps fully qualified type paths to their package information
	Packages []*PackageWithTemplates
	// Error encountered during type resolution, if any
	Err error
}

// NewRegistry creates a new Registry
func NewRegistry(pkgWithTemplatesList []*PackageWithTemplates) *Registry {
	return &Registry{
		Packages: pkgWithTemplatesList,
	}
}

func NewEmptyRegistry() *Registry {
	return &Registry{
		Packages: []*PackageWithTemplates{},
	}
}

func (r *Registry) AddPackage(pkg *PackageWithTemplates) {
	r.Packages = append(r.Packages, pkg)
}

type InMemoryPackageOpts struct {
	PackagePath   string
	PackageName   string
	TemplateFiles map[string]string
	Types         []*types.TypeName
}

func (r *PackageWithTemplates) MustAddAndParseTemplates(ctx context.Context, files map[string]string) {
	err := r.AddAndParseTemplates(ctx, files)
	if err != nil {
		panic(err)
	}
}

func (r *PackageWithTemplates) AddAndParseTemplates(ctx context.Context, files map[string]string) error {
	for name, file := range files {
		tmpl, err := ParseTemplate(ctx, name, file)
		if err != nil {
			return errors.Errorf("parsing in memory template %s: %w", name, err)
		}
		r.Templates[name] = tmpl
	}
	return nil
}

func (r *PackageWithTemplates) AddTypes(types []*types.TypeName) {
	for _, obj := range types {
		r.Package.Types.Scope().Insert(obj)
	}
}

func (r *PackageWithTemplates) AddStruct(name string, fieldMap map[string]types.Type) *types.Named {
	fields := make([]*types.Var, 0, len(fieldMap))
	for name, fieldType := range fieldMap {
		fields = append(fields, types.NewField(0, r.Package.Types, name, fieldType, false))
	}

	structed := types.NewStruct(fields, nil)

	named := types.NewNamed(
		types.NewTypeName(0, r.Package.Types, name, nil),
		structed,
		nil,
	)

	r.AddTypes([]*types.TypeName{named.Obj()})

	return named
}

func (r *Registry) AddInMemoryPackageForTesting(ctx context.Context, path string) *PackageWithTemplates {
	name := filepath.Base(path)
	pkg := packages.Package{
		PkgPath: path,
		Name:    name,
	}

	pkg.Types = types.NewPackage(path, name)

	pkgWithTemplates := &PackageWithTemplates{
		Package:   &pkg,
		Templates: map[string]*template.Template{},
	}

	r.Packages = append(r.Packages, pkgWithTemplates)

	return pkgWithTemplates
}

// GetPackage returns a package by name
func (r *Registry) GetPackage(ctx context.Context, packageName string) (*types.Package, error) {
	// zerolog.Ctx(ctx).Debug().Str("packageName", packageName).Interface("packages", r.Packages).Msg("looking for package")

	// First, try to find an exact match
	for _, pkg := range r.Packages {
		if pkg.Package.PkgPath == packageName {
			zerolog.Ctx(ctx).Debug().Str("package", packageName).Msg("found exact match")
			return pkg.Package.Types, nil
		}
	}

	// Try to find by package name
	for _, pkg := range r.Packages {
		if path.Base(pkg.Package.PkgPath) == packageName {
			zerolog.Ctx(ctx).Debug().Str("packageName", packageName).Str("path", pkg.Package.PkgPath).Msg("found by name")
			return pkg.Package.Types, nil
		}
	}

	// Try to find by path suffix
	for _, pkg := range r.Packages {
		if strings.HasSuffix(pkg.Package.PkgPath, "/"+packageName) {
			zerolog.Ctx(ctx).Debug().Str("packageName", packageName).Str("path", pkg.Package.PkgPath).Msg("found by suffix")
			return pkg.Package.Types, nil
		}
	}

	zerolog.Ctx(ctx).Debug().Str("packageName", packageName).Msg("not found")
	return nil, errors.Errorf("package %s not found", packageName)
}

// GetTypes retrieves all types from a package
func (r *Registry) GetTypes(ctx context.Context, pkgPath string) (map[string]types.Object, error) {
	pkg, err := r.GetPackage(ctx, pkgPath)
	if err != nil {
		return nil, err
	}

	types := make(map[string]types.Object)
	scope := pkg.Scope()
	for _, name := range scope.Names() {
		types[name] = scope.Lookup(name)
	}

	return types, nil
}

// TypeExists checks if a type exists in the registry
func (r *Registry) TypeExists(typePath string) bool {
	for _, pkg := range r.Packages {
		path := pkg.Package.Types.Path()
		for _, name := range pkg.Package.Types.Scope().Names() {
			if path+"."+name == typePath {
				return true
			}
		}
	}

	return false
}

// GetFieldType returns the type of a field in a struct type
func (r *Registry) GetFieldType(structType *types.Struct, fieldName string) (types.Type, error) {
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		if field.Name() == fieldName {
			return field.Type(), nil
		}
	}
	return nil, errors.Errorf("field %s not found", fieldName)
}
