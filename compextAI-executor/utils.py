from pydantic import BaseModel, Field, create_model
from typing import Dict, Any, Type, Union, Optional, List

def create_pydantic_model_from_dict(
    model_name: str,
    schema_dict: Dict[str, Any]
) -> Type[BaseModel]:
    """
    Recursively creates a Pydantic BaseModel from a dictionary representing the JSON schema.
    """
    def parse_schema(name: str, schema: Dict[str, Any]) -> Type[Any]:
        json_type = schema.get('type', 'string')

        if json_type == 'object':
            properties = schema.get('properties', {})
            required_fields = schema.get('required', [])
            fields = {}
            for field_name, field_info in properties.items():
                field_type = parse_schema(f"{name}_{field_name}", field_info)
                is_required = field_name in required_fields
                description = field_info.get('description', '')
                default = field_info.get('default', ...)
                if not is_required:
                    field_type = Optional[field_type]
                    default = None if default is ... else default
                fields[field_name] = (field_type, Field(default, description=description))
            model = create_model(name, **fields)
            return model

        elif json_type == 'array':
            items_schema = schema.get('items', {})
            items_type = parse_schema(f"{name}_Item", items_schema)
            return List[items_type]

        else:
            type_mapping = {
                'string': str,
                'number': float,
                'integer': int,
                'boolean': bool,
                'null': None,
            }
            python_type = type_mapping.get(json_type, Any)
            return python_type

    model = parse_schema(model_name, schema_dict)
    return model
