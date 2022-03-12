import os
import tempfile
import unittest

import olca_schema.zipio as zipio
from olca_schema import *


class ZipioTest(unittest.TestCase):

    def test_simple(self):
        zip_file = tempfile.mktemp('.zip')
        actor = actor_of('My company')
        print(zip_file)
        with zipio.ZipWriter(zip_file) as writer:
            writer.write(actor)
        with zipio.ZipReader(zip_file) as reader:
            clone = reader.read_actor(actor.id)
            self.assertEqual(actor.id, clone.id)
            self.assertEqual('My company', clone.name)
            none = reader.read_actor('does not exist')
            self.assertIsNone(none)
        os.remove(zip_file)


if __name__ == '__main__':
    unittest.main()
