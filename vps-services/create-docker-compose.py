import yaml
import os

def save_yaml(data, file_path):
    with open(file_path, 'w') as file:
        yaml.safe_dump(data, file, default_flow_style=False)

def modify_docker_compose(data, changes):
    for service, values in changes.items():
        if service in data['services']:
            data['services'][service].update(values)
        else:
            data['services'][service] = values
    return data

file_path = 'docker-compose.yml'

service_name = os.getenv('SERVICE_NAME')
meilisearch_image_url = os.getenv('MEILISEARCH_IMAGE_URL')
meilisearch_host = os.getenv('MEILISEARCH_HOST')
meilisearch_master_key = os.getenv('MEILISEARCH_MASTER_KEY')
meilisearch_no_analytics = os.getenv('MEILISEARCH_NO_ANALYTICS')
meilisearch_env = os.getenv('MEILISEARCH_ENV')

if not service_name:
    raise ValueError('SERVICE_NAME is required')
if not meilisearch_image_url:
    raise ValueError('MEILISEARCH_IMAGE_URL is required')
if not meilisearch_host:
    raise ValueError('MEILISEARCH_HOST is required')
if not meilisearch_master_key:
    raise ValueError('MEILISEARCH_MASTER_KEY is required')
if not meilisearch_no_analytics:
    raise ValueError('MEILISEARCH_NO_ANALYTICS is required')
if not meilisearch_env:
    raise ValueError('MEILISEARCH_ENV is required')

changes = {
    'meilisearch': {
        'container_name': service_name,
        'image': meilisearch_image_url,
        'environment': {
            'MEILI_HOST': meilisearch_host,
            'MEILI_MASTER_KEY': meilisearch_master_key,
            'MEILI_NO_ANALYTICS': meilisearch_no_analytics,
            'MEILI_ENV': meilisearch_env,
            'MEILI_LOG_LEVEL': 'INFO',
            'MEILI_DB_PATH': '/data.ms',
            'TZ': 'America/Sao_Paulo'
        },
        'ports': [
            '7700:7700'
        ],
        'networks': [
            'meilisearch'
        ],
        'volumes': [
            './data.ms:/data.ms'
        ],
        'restart': 'unless-stopped'
    }
}

data = {
    "services": {
    },
    "networks": {
        "meilisearch": {
            "driver": "bridge"
        }
    }
}

modified_data = modify_docker_compose(data, changes)

save_yaml(modified_data, file_path)

print("docker-compose.yml file updated successfully")
