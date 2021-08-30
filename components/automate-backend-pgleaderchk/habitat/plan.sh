# shellcheck disable=SC2148
pkg_name=automate-backend-pgleaderchk
binary_name=pgleaderchk
pkg_origin=chef
pkg_repo=develop
pkg_version="0.1.2"
pkg_description="Automate Backend PostreSQL leader check"
pkg_maintainer="Chef Software Inc. <support@chef.io>"
pkg_license=('Chef-MLSA')
pkg_upstream_url="http://github.com/punitmundra/develop/components/automate-backend-pgleaderchk"
pkg_deps=(
  core/bash
  core/curl
  core/jq-static
  chef/mlsa
)
pkg_build_deps=(
  core/gcc
)
pkg_exports=(
  [port]=httpd.port
)
pkg_binds_optional=(
  [database]="port"
  [database]="ssl"
)
pkg_exposes=(port)


pkg_bin_dirs=(bin)
pkg_scaffolding="${local_scaffolding_origin:-chef}/automate-scaffolding-go"
scaffolding_go_base_path=github.com/punitmundra
scaffolding_go_repo_name=develop
scaffolding_go_import_path="${scaffolding_go_base_path}/${scaffolding_go_repo_name}/components/${pkg_name}"
scaffolding_go_binary_list=(
  "${scaffolding_go_import_path}/cmd/${pkg_name}"
)

do_install() {
  do_default_install
  build_line "Overriding Install process"
  cp -r "${GOBIN}/$binary_name" "${pkg_prefix}/bin"
}

do_strip() {
  return 0
}
