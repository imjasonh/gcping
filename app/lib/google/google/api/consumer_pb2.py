# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: google/api/consumer.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
from google.protobuf import descriptor_pb2
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='google/api/consumer.proto',
  package='google.api',
  syntax='proto3',
  serialized_pb=_b('\n\x19google/api/consumer.proto\x12\ngoogle.api\"=\n\x11ProjectProperties\x12(\n\nproperties\x18\x01 \x03(\x0b\x32\x14.google.api.Property\"\xac\x01\n\x08Property\x12\x0c\n\x04name\x18\x01 \x01(\t\x12/\n\x04type\x18\x02 \x01(\x0e\x32!.google.api.Property.PropertyType\x12\x13\n\x0b\x64\x65scription\x18\x03 \x01(\t\"L\n\x0cPropertyType\x12\x0f\n\x0bUNSPECIFIED\x10\x00\x12\t\n\x05INT64\x10\x01\x12\x08\n\x04\x42OOL\x10\x02\x12\n\n\x06STRING\x10\x03\x12\n\n\x06\x44OUBLE\x10\x04\x42!\n\x0e\x63om.google.apiB\rConsumerProtoP\x01\x62\x06proto3')
)
_sym_db.RegisterFileDescriptor(DESCRIPTOR)



_PROPERTY_PROPERTYTYPE = _descriptor.EnumDescriptor(
  name='PropertyType',
  full_name='google.api.Property.PropertyType',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='UNSPECIFIED', index=0, number=0,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='INT64', index=1, number=1,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='BOOL', index=2, number=2,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='STRING', index=3, number=3,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='DOUBLE', index=4, number=4,
      options=None,
      type=None),
  ],
  containing_type=None,
  options=None,
  serialized_start=201,
  serialized_end=277,
)
_sym_db.RegisterEnumDescriptor(_PROPERTY_PROPERTYTYPE)


_PROJECTPROPERTIES = _descriptor.Descriptor(
  name='ProjectProperties',
  full_name='google.api.ProjectProperties',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='properties', full_name='google.api.ProjectProperties.properties', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=41,
  serialized_end=102,
)


_PROPERTY = _descriptor.Descriptor(
  name='Property',
  full_name='google.api.Property',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='google.api.Property.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='type', full_name='google.api.Property.type', index=1,
      number=2, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='description', full_name='google.api.Property.description', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
    _PROPERTY_PROPERTYTYPE,
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=105,
  serialized_end=277,
)

_PROJECTPROPERTIES.fields_by_name['properties'].message_type = _PROPERTY
_PROPERTY.fields_by_name['type'].enum_type = _PROPERTY_PROPERTYTYPE
_PROPERTY_PROPERTYTYPE.containing_type = _PROPERTY
DESCRIPTOR.message_types_by_name['ProjectProperties'] = _PROJECTPROPERTIES
DESCRIPTOR.message_types_by_name['Property'] = _PROPERTY

ProjectProperties = _reflection.GeneratedProtocolMessageType('ProjectProperties', (_message.Message,), dict(
  DESCRIPTOR = _PROJECTPROPERTIES,
  __module__ = 'google.api.consumer_pb2'
  # @@protoc_insertion_point(class_scope:google.api.ProjectProperties)
  ))
_sym_db.RegisterMessage(ProjectProperties)

Property = _reflection.GeneratedProtocolMessageType('Property', (_message.Message,), dict(
  DESCRIPTOR = _PROPERTY,
  __module__ = 'google.api.consumer_pb2'
  # @@protoc_insertion_point(class_scope:google.api.Property)
  ))
_sym_db.RegisterMessage(Property)


DESCRIPTOR.has_options = True
DESCRIPTOR._options = _descriptor._ParseOptions(descriptor_pb2.FileOptions(), _b('\n\016com.google.apiB\rConsumerProtoP\001'))
# @@protoc_insertion_point(module_scope)
