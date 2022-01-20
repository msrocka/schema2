@echo off

if exist docs (
    rmdir /S /Q docs
)

py build_from_yaml.py
mkdir docs\html
generate-schema-doc docs docs/html --config with_footer=false
