import json
import os
from typing import NamedTuple, Optional
import yaml


YAML_DIR = 'C:/Users/Win10/Projects/openLCA/repos/olca-schema/yaml'


class Model(NamedTuple):
    classes: dict[str, 'ClassDef']
    enums: dict[str, 'EnumDef']

    @staticmethod
    def new() -> 'Model':
        return Model({}, {})


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
                case _:
                    if type_def[0].isupper():
                        return {'$ref': f'{type_def}.schema.json'}
                    print(f'unmatched primitive type: {type_def}')
                    return type_def

        schema = {}
        type_def: str = self.obj['type']
        if type_def.startswith('Ref['):
            schema['type'] = {
                '$ref': 'Ref.schema.json'
            }
        elif type_def.startswith('List['):
            item_type = type_def[5:len(type_def)-1]
            schema['type'] = 'array'
            schema['item'] = schema_type(item_type)
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


class EnumDef:
    model: 'Model'
    pass


def main():
    model = Model.new()
    for f in os.listdir(YAML_DIR):
        path = os.path.join(YAML_DIR, f)
        with open(path, 'r', encoding='utf-8') as inp:
            decl: dict[str, any] = yaml.load(inp, yaml.SafeLoader)
            class_def: Optional[dict] = decl.get('class')
            if not class_def:
                continue
            name: str = class_def['name']
            model.classes[name] = ClassDef(model, class_def)

    for class_def in model.classes.values():
        name = class_def.name()
        schema = class_def.to_schema()
        path = f'generated/{name}.schema.json'
        with open(path, 'w', encoding='utf-8') as out:
            json.dump(schema, out, indent=2)


if __name__ == '__main__':
    main()
