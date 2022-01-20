# you need to have PyYaml installed in order to run this script

import json
import os
import yaml

from typing import NamedTuple, Optional


YAML_DIR = 'C:/Users/Win10/Projects/openLCA/repos/olca-schema/yaml'
OUT_DIR = 'docs'

class Model(NamedTuple):
    classes: dict[str, 'ClassDef']
    enums: dict[str, 'EnumDef']

    @staticmethod
    def new() -> 'Model':
        return Model({}, {})

    def all_defs(self) -> list:
        defs = []
        for c in self.classes.values():
            defs.append(c)
        for e in self.enums.values():
            defs.append(e)
        return defs


class PropDef(NamedTuple):
    model: Model
    class_def: 'ClassDef'
    obj: dict[str, any]

    def name(self) -> str:
        return self.obj['name']

    def to_schema(self) -> dict:

        def schema_type(type_def: str) -> str | dict:
            match type_def:
                case 'string': return 'string'
                case 'double': return 'number'
                case 'boolean': return 'boolean'
                case 'int': return 'integer'
                case 'integer': return 'integer'
                case _:
                    if type_def.startswith('Ref['):
                        return {'$ref': f'Ref.schema.json'}
                    if type_def[0].isupper():
                        return {'$ref': f'{type_def}.schema.json'}
                    print(f'unmatched primitive type: {type_def}')
                    return type_def

        schema = {}
        type_def: str = self.obj['type']
        if type_def.startswith('Ref['):
            schema['$ref'] = 'Ref.schema.json'
        elif type_def.startswith('List['):
            item_type = type_def[5:len(type_def)-1]
            schema['type'] = 'array'
            schema['item'] = schema_type(item_type)
        elif type_def == 'dateTime':
            schema['type'] = 'string'
            schema['format'] = 'date-time'
        elif type_def == 'date':
            schema['type'] = 'string'
            schema['format'] = 'date'
        elif type_def == 'GeoJSON':
            schema['type'] = 'object'
            schema['format'] = 'GeoJSON'
        elif type_def[0].isupper():
            schema['$ref'] = f'{type_def}.schema.json'
        else:
            schema['type'] = schema_type(type_def)

        doc = self.obj.get('doc')
        if doc:
            schema['description'] = doc

        return schema


class ClassDef(NamedTuple):
    model: Model
    obj: dict[str, any]

    def name(self) -> str:
        return self.obj['name']

    def doc(self) -> Optional[str]:
        return self.obj.get('doc')

    def super_class(self) -> Optional['ClassDef']:
        c: str = self.obj.get('superClass')
        if not c:
            return None
        return self.model.classes.get(c)

    def properties(self) -> list[PropDef]:
        sc = self.super_class()
        properties = sc.properties() if sc else []
        props: list[dict[str, any]] = self.obj.get('properties')
        if not props:
            return properties
        for prop in props:
            properties.append(PropDef(self.model, self, prop))
        return properties

    def to_schema(self) -> dict:
        name = self.name()
        schema = {
            '$id': f'{name}.schema.json',
            '$schema': 'https://json-schema.org/draft/2020-12/schema',
            'type': 'object',
            'title': name,
            'properties': {}
        }

        doc = self.doc()
        if doc:
            schema['description'] = doc

        proto_idx = 0
        for prop in self.properties():
            proto_idx += 1
            prop_schema = prop.to_schema()
            prop_schema['protoIndex'] = proto_idx
            schema['properties'][prop.name()] = prop_schema

        return schema


class EnumDef(NamedTuple):
    obj: dict[str, any]

    def name(self):
        return self.obj['name']

    def to_schema(self) -> dict:
        name = self.name()
        schema = {
            '$id': f'{name}.schema.json',
            '$schema': 'https://json-schema.org/draft/2020-12/schema',
            'type': 'string',
            'title': name,
            'enum': []
        }

        doc = self.obj.get('doc')
        if doc:
            schema['description'] = doc

        items: list[dict[str, any]] = self.obj.get('items', [])
        for item in items:
            schema['enum'].append(item['name'])

        return schema


def main():
    if not os.path.isdir(OUT_DIR):
        os.makedirs(OUT_DIR)

    model = Model.new()
    for f in os.listdir(YAML_DIR):
        path = os.path.join(YAML_DIR, f)
        with open(path, 'r', encoding='utf-8') as inp:
            decl: dict[str, any] = yaml.load(inp, yaml.SafeLoader)
            class_def: Optional[dict] = decl.get('class')
            if class_def:
                name: str = class_def['name']
                model.classes[name] = ClassDef(model, class_def)
            enum_def: Optional[dict] = decl.get('enum')
            if enum_def:
                name: str = enum_def['name']
                model.enums[name] = EnumDef(enum_def)

    for d in model.all_defs():
        name = d.name()
        schema = d.to_schema()
        path = f'{OUT_DIR}/{name}.schema.json'
        with open(path, 'w', encoding='utf-8') as out:
            json.dump(schema, out, indent=2)


if __name__ == '__main__':
    main()
