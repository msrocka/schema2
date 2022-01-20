import json
import os

if __name__ == '__main__':
    root_types: list[str] = []
    for f in os.listdir('docs'):
        if not f.endswith('.json'):
            continue
        with open(f'docs/{f}', 'r', encoding='utf-8') as inp:
            obj: dict[str, any] = json.load(inp)
            props: dict = obj.get('properties')
            if not props:
                continue
            if '@id' in props:
                root_types.append(obj['title'])

    base_url = 'https://msrocka.github.io/schema2'
    for t in root_types:
        print(f'* {t}: [schema]({base_url}/{t}.schema.json) [doc]({base_url}/html/{t}.schema.html)')
