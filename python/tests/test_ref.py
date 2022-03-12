import unittest

from olca_schema import *


class RefTest(unittest.TestCase):

    def setUp(self):
        units = unit_group_of('Units of mass', 'kg')
        mass = flow_property_of('Mass', units)
        flow = elementary_flow_of(name="CO2", flow_property=mass)
        flow.category = 'air/unspefified'
        self.flow = flow

    def test_to_ref(self):
        ref = self.flow.to_ref()
        self.assertEqual('Flow', ref.model_type)
        self.assertEqual(self.flow.id, ref.id)
        self.assertEqual(self.flow.name, ref.name)
        self.assertEqual(self.flow.category, ref.category)

    def test_to_dict(self):
        ref_dict = self.flow.to_ref().to_dict()
        self.assertEqual('Flow', ref_dict['@type'])
        self.assertEqual(self.flow.id, ref_dict['@id'])
        self.assertEqual(self.flow.name, ref_dict['name'])
        self.assertEqual(self.flow.category, ref_dict['category'])

    def test_from_dict(self):
        ref_dict = self.flow.to_ref().to_dict()
        ref = Ref.from_dict(ref_dict)
        self.assertEqual('Flow', ref.model_type)
        self.assertEqual(self.flow.id, ref.id)
        self.assertEqual(self.flow.name, ref.name)
        self.assertEqual(self.flow.category, ref.category)


if __name__ == '__main__':
    unittest.main()
