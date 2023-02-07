/*
Copyright 2023 The KubeService-Stack Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mcpack

const (
	MCPACKV2_INVALID      = 0x00
	MCPACKV2_OBJECT       = 0x10
	MCPACKV2_ARRAY        = 0x20
	MCPACKV2_STRING       = 0x50
	MCPACKV2_BINARY       = 0x60
	MCPACKV2_INT8         = 0x11
	MCPACKV2_INT16        = 0x12
	MCPACKV2_INT32        = 0x14
	MCPACKV2_INT64        = 0x18
	MCPACKV2_UINT8        = 0x21
	MCPACKV2_UINT16       = 0x22
	MCPACKV2_UINT32       = 0x24
	MCPACKV2_UINT64       = 0x28
	MCPACKV2_BOOL         = 0x31
	MCPACKV2_FLOAT        = 0x44
	MCPACKV2_DOUBLE       = 0x48
	MCPACKV2_DATE         = 0x58
	MCPACKV2_NULL         = 0x61
	MCPACKV2_SHORT_ITEM   = 0x80
	MCPACKV2_FIXED_ITEM   = 0xf0
	MCPACKV2_DELETED_ITEM = 0x70

	MCPACKV2_SHORT_STRING = MCPACKV2_STRING | MCPACKV2_SHORT_ITEM
	MCPACKV2_SHORT_BINARY = MCPACKV2_BINARY | MCPACKV2_SHORT_ITEM

	MCPACKV2_KEY_MAX_LEN = 254

	MAX_SHORT_VITEM_LEN = 255
)
