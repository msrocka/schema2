import json
import zipfile

import olca_schema as schema


class ZipWriter:

    def __init__(self, file_name: str):
        self.__zip = zipfile.ZipFile(
            file_name, mode='a', compression=zipfile.ZIP_DEFLATED)

    def __enter__(self):
        return self

    def __exit__(self, type, value, traceback):
        self.close()

    def close(self):
        self.__zip.close()

    def write(self, entity: schema.RootEntity):
        if entity.id is None or entity.id == '':
            raise ValueError('entity must have an ID')
        folder = _folder_of_entity(entity)
        path = f'{folder}/{entity.id}.json'
        data = json.dumps(entity.to_dict(), indent='  ')
        self.__zip.writestr(path, data)


def _folder_of_entity(entity: schema.RootEntity):
    if entity is None:
        raise ValueError("unknown root entity type: " + entity)
    t = type(entity)
    if t == schema.Actor:
        return 'actors'
    if t == schema.Currency:
        return 'currencies'
    if t == schema.DQSystem:
        return 'dq_systems'
    if t == schema.Epd:
        return 'epds'
    if t == schema.Flow:
        return 'flows'
    if t == schema.FlowProperty:
        return 'flow_properties'
    if t == schema.ImpactCategory:
        return 'impact_categories'
    if t == schema.ImpactMethod:
        return 'impact_methods'
    if t == schema.Location:
        return 'locations'
    if t == schema.Parameter:
        return 'parameters'
    if t == schema.Process:
        return 'processes'
    if t == schema.ProductSystem:
        return 'product_systems'
    if t == schema.Project:
        return 'projects'
    if t == schema.Result:
        return 'results'
    if t == schema.SocialIndicator:
        return 'social_indicators'
    if t == schema.Source:
        return 'sources'
    if t == schema.UnitGroup:
        return 'unit_groups'
    raise ValueError("unknown entity type %s" % t)
