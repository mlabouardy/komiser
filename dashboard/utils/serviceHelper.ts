// Define a union type for supported cloud providers.
export type Providers =
  | 'aws'
  | 'gcp'
  | 'digitalocean'
  | 'azure'
  | 'civo'
  | 'kubernetes'
  | 'linode'
  | 'tencent'
  | 'oci'
  | 'scaleway'
  | 'mongodbatlas'
  | 'ovh'
  | 'scaleway'
  | 'tencent';

// Define an object that maps each provider to an array of its supported services.
export const allProvidersServices: { [key in Providers]: string[] } = {
  aws: [
    'api gateway',
    'cloudfront',
    'cloudwatch',
    'cloudwatch dashboard',
    'cloudwatch log group',
    'cloudwatch metric stream',
    'codebuild',
    'codecommit',
    'codedeploy',
    'dynamodb',
    'acl',
    'elastic ip',
    'ec2',
    'internet gateway',
    'ec2 keypair',
    'network interface',
    'ec2 placement group',
    'security group',
    'ec2 snapshot',
    'ec2 spot instance request',
    'subnet',
    'ebs',
    'vpc',
    'vpc endpoint',
    'vpc peering connection',
    'ecr',
    'ecs cluster',
    'ecs container instance',
    'ecs task definition',
    'efs',
    'eks',
    'elasticache',
    'target group',
    'iam group',
    'iam instance profile',
    'iam identity provider',
    'iam policy',
    'iam role',
    'iam saml provider',
    'kinesis stream',
    'kinesis efo consumer',
    'kms',
    'eventsourcemapping',
    'lambda',
    'opensearch service domain',
    'rds backup',
    'rds',
    'rds cluster snapshot',
    'rds db proxy endpoint',
    'rds instance',
    'rds proxy',
    'rds snapshot',
    'redshift eventsubscription',
    's3',
    'service catalog application',
    'sns',
    'sqs',
    'ssm association',
    'ssm maintenance window'
  ],
  azure: [
    'disk',
    'image',
    'snapshot',
    'virtual machine',
    'sql database server',
    'application gateway',
    'firewall',
    'load balancer',
    'databox',
    'queue',
    'local network gateway'
  ],
  civo: [
    'compute',
    'kubernetes',
    'firewall',
    'load balancer',
    'database',
    'diskimage',
    'object store',
    'volume',
    'network'
  ],
  digitalocean: [
    'database',
    'droplet',
    'namespace',
    'trigger',
    'kubernetes',
    'kubernetes cluster',
    'firewall',
    'load balancer',
    'vpc',
    'volume'
  ],
  gcp: [
    'bigquery table',
    'certificate',
    'compute disk',
    'vm instance',
    'compute disk snapshot',
    'cluster',
    'cloud functions',
    'api gateways',
    'iam roles',
    'iam service accounts',
    'kms key',
    'redis',
    'sql instance',
    'bucket'
  ],
  kubernetes: [
    'daemonset',
    'deployment',
    'ingress',
    'job',
    'namespace',
    'node',
    'pod',
    'persistentvolume',
    'persistentvolumeclaim',
    'serviceaccount',
    'service',
    'statefulset'
  ],
  linode: [
    'linode instance',
    'lke',
    'firewall',
    'nodebalancer',
    'postgresql',
    'mysql',
    'bucket',
    'database',
    'volume'
  ],
  mongodbatlas: ['cluster', 'serverless cluster'],
  oci: [
    'vm',
    'application',
    'function',
    'identity policy',
    'autonomous database',
    'block volume',
    'objectstorage bucket'
  ],
  ovh: [
    'alerting',
    'image',
    'instance',
    'kube',
    'ip',
    'network',
    'vrack',
    'project',
    'ssh',
    'ssl',
    'container',
    'volume',
    'user'
  ],
  scaleway: [
    'server',
    'kubernetes',
    'containerregistry',
    'loadbalancer',
    'serverlesscontainer',
    'function',
    'database'
  ],
  tencent: ['instance']
};

/**
 * Check if a specific service is supported by a given provider.
 * @param {Providers} provider - The cloud provider.
 * @param {string} service - The service to check for support.
 * @returns {boolean} - Returns true if the service is supported, otherwise returns false.
 *
 * This function converts the provider and service names to lowercase to ensure case-insensitive matching.
 * It then looks up the array of supported services for the given provider and checks if the specified service is included.
 * The result is a boolean indicating whether the service is supported.
 */
export function checkIfServiceIsSupported(
  provider: string,
  service: string
): boolean {
  const lowercaseProvider = provider.toLowerCase() as Providers;
  const lowercaseService = service.toLowerCase();
  const services = allProvidersServices[lowercaseProvider];
  return services.includes(lowercaseService);
}

/**
 * Checks if there are any services in the receivedServices array that are not supported.
 * @param {string[]} receivedServices - An array of services to check for support.
 * @returns {boolean} - Returns true if at least one service is not supported, otherwise returns false.
 *
 * This function compares the receivedServices against a list of all supported services
 * (flattened from the allProvidersServices object). It returns true if any service in
 * receivedServices is not found in the list of supported services.
 */
export function checkIsSomeServiceUnavailable(
  receivedServices: string[]
): boolean {
  const allServices = Object.values(allProvidersServices).flat();
  return receivedServices.some(
    service => !allServices.includes(service.toLowerCase())
  );
}
