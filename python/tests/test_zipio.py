import os
import tempfile
import unittest

import olca_schema.zipio as zipio
from olca_schema import *


class ZipioTest(unittest.TestCase):

    def test_simple(self):
        zip_file = tempfile.mktemp('.zip')
        actor = actor_of('My company')
        with zipio.ZipWriter(zip_file) as writer:
            writer.write(actor)
        with zipio.ZipReader(zip_file) as reader:
            clone = reader.read_actor(actor.id)
            self.assertEqual(actor.id, clone.id)
            self.assertEqual('My company', clone.name)
            none = reader.read_actor('does not exist')
            self.assertIsNone(none)
        os.remove(zip_file)

    def test_io(self):
        zip_file = tempfile.mktemp('.zip')
        classes = [
            schema.Actor, schema.Currency, schema.DQSystem, schema.Epd,
            schema.Flow, schema.FlowProperty, schema.ImpactCategory,
            schema.ImpactMethod, schema.Location, schema.Parameter,
            schema.Process, schema.ProductSystem, schema.Project,
            schema.Result, schema.SocialIndicator, schema.Source,
            schema.UnitGroup]

        def uid(c): return c.__name__.lower()
        with zipio.ZipWriter(zip_file) as writer:
            for c in classes:
                instance = c(id=uid(c), name=uid(c))
                writer.write(instance)
        with zipio.ZipReader(zip_file) as reader:
            for c in classes:
                instance = reader.read(c, uid(c))
                self.assertEqual(c, type(instance))
                self.assertEqual(uid(c), instance.id)
                self.assertEqual(uid(c), instance.name)
        os.remove(zip_file)


if __name__ == '__main__':
    unittest.main()
