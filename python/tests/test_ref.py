import unittest

from olca_schema import *


class RefTest(unittest.TestCase):

    def test_to_ref(self):
        units = unit_group_of('Units of mass', 'kg')
        mass = flow_property_of('Mass', units)
        flow = elementary_flow_of(name="CO2", flow_property=mass)
        ref = flow.to_ref()
        self.assertEqual(flow.name, ref.name)
        self.assertEqual(flow.id, ref.id)
        self.assertEqual('Flow', ref.model_type)
        ref_dict = ref.to_dict()
        self.assertEqual(flow.name, ref_dict['name'])
        self.assertEqual(flow.id, ref_dict['@id'])
        self.assertEqual('Flow', ref_dict['@type'])


if __name__ == '__main__':
    unittest.main()
