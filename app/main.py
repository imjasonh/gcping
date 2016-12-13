#!/usr/bin/env python

from google.cloud import bigquery
from google.cloud.bigquery.schema import SchemaField

import json
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
  SchemaField('user_agent', 'STRING',    'NULLABLE', 'Browser user-agent making ping requests'),
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
    ua = self.request.headers.get('User-Agent')
    data = json.loads(self.request.body)
    rows = []
    for i in data:
      rows.append((
        i.timestamp,
        i.took,
        data.lat,
        data.lng,
        i.region,
        ua,
      ))
    table.insert_rows(rows, skip_invalid_rows=True, ignore_unknown_values=True)

app = webapp2.WSGIApplication([('/api', API)])
