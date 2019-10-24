#!/usr/bin/env ruby

require 'tmpdir'
require 'fileutils'

release_path=ARGV[0]
release_filename=File.basename(release_path, '.tgz')
release_name=release_filename.split('-').first
release_version=/(\d+\.\d+\.\d+)/.match(release_filename)[1]

warn "release_filename: #{release_filename}"
warn "release_name: #{release_name}"
warn "release_version: #{release_version}"

manifest = `bash -c 'go run main.go --path #{release_path} --merge <(bosh int example/metadata.yml -v bosh_release_name=#{release_name} -v bosh_release_filename=#{release_filename}.tgz -v bosh_release_version=#{release_version})'`

dir = Dir.mktmpdir
metadata_dir = File.join(dir, 'metadata')
FileUtils.mkdir_p(metadata_dir)
File.write(File.join(metadata_dir, 'metadata.yml'), manifest)

File.write('metadata.yml', manifest)

releases_dir = File.join(dir, 'releases')
FileUtils.mkdir_p(releases_dir)
FileUtils.cp(release_path, File.join(releases_dir, File.basename(release_path)))

migrations_dir = File.join(dir, 'migrations')
FileUtils.mkdir_p(File.join(migrations_dir, 'v1'))

workspace = Dir.pwd
product_path = File.join(workspace, 'example-0.0-build.0.pivotal')
FileUtils.rm(product_path) rescue nil
Dir.chdir(dir) do
    system("zip -r #{product_path} migrations releases metadata")
end