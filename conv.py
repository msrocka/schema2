import json
import os
from typing import Optional
import yaml


YAML_DIR = 'C:/Users/Win10/Projects/openLCA/repos/olca-schema/yaml'

for f in os.listdir(YAML_DIR):
    path = os.path.join(YAML_DIR, f)
    with open(path, 'r', encoding='utf-8') as stream:
        decl: dict[str, any] = yaml.load(stream, yaml.SafeLoader)
        class_def: Optional[dict] = decl.get('class')
        if not class_def:
            continue
        name: str = class_def['name']
        class_schema = {
            '$id': f'{name}.schema.json',
            '$schema': 'https://json-schema.org/draft/2020-12/schema',
            'properties': {}
        }

        doc: str = class_def.get('doc')
        if doc:
            class_schema['description'] = doc

        super_class: str = class_def.get('superClass')
        if super_class and super_class != 'Entity':
            class_schema['allOf'] = [{'$ref': f'{super_class}.schema.json'}]

        outpath = f'generated/{name}.schema.json'
        with open(outpath, 'w', encoding='utf-8') as out:
            json.dump(class_schema, out, indent=2)

        print(name)

