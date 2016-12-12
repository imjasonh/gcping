#!/usr/bin/env python

from google.cloud import bigquery
from google.cloud.bigquery.schema import SchemaField

import logging
import webapp2


_DATASET = 'gcping'
_TABLE = 'results'
_SCHEMA = [
  SchemaField('timestamp',  'TIMESTAMP', 'NULLABLE', 'Time the ping request finished'),
  SchemaField('latency_ms', 'INTEGER',   'NULLABLE', 'Amount of time in milliseconds the ping took'),
  SchemaField('lat',        'FLOAT',     'NULLABLE', 'Latitude of the user making the ping request'),
  SchemaField('lng',        'FLOAT',     'NULLABLE', 'Longitude of the user making the ping request'),
  SchemaField('region',     'STRING',    'NULLABLE', 'GCP region being pinged'),
]


class API(webapp2.RequestHandler):
  def post(self):
    client = bigquery.Client()
    dataset = client.dataset(_DATASET)
    if not dataset.exists():
      dataset.create()
      logging.info('Dataset %s created', dataset.name)
    table = dataset.table(_TABLE, _SCHEMA)
    if not table.exists():
      table.create()
      logging.info('Table %s created', table.name)
    self.response.out.write('Dataset and table exist')

app = webapp2.WSGIApplication([('/api', API)])
