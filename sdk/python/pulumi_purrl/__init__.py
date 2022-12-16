# coding=utf-8
# *** WARNING: this file was generated by pulumi. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

from . import _utilities
import typing
# Export this package's modules as members:
from .provider import *

# Make subpackages available:
if typing.TYPE_CHECKING:
    import pulumi_purrl.purrl as __purrl
    purrl = __purrl
else:
    purrl = _utilities.lazy_import('pulumi_purrl.purrl')

_utilities.register(
    resource_modules="""
[
 {
  "pkg": "purrl",
  "mod": "purrl",
  "fqn": "pulumi_purrl.purrl",
  "classes": {
   "purrl:purrl:Purrl": "Purrl"
  }
 }
]
""",
    resource_packages="""
[
 {
  "pkg": "purrl",
  "token": "pulumi:providers:purrl",
  "fqn": "pulumi_purrl",
  "class": "Provider"
 }
]
"""
)
