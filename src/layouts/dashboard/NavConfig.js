// component
import Iconify from '../../components/Iconify';

// ----------------------------------------------------------------------

const getIcon = (name) => <Iconify icon={name} width={22} height={22} />;

const navConfig = [
  {
    title: 'Warm Function Invocations',
    path: '/dashboard/warm/aws',
    icon: getIcon('bxs:hot'),
  },
  {
    title: 'Cold Function Invocations',
    path: '/dashboard/cold/',
    icon: getIcon('ic:baseline-severe-cold'),
    children:[
      {
        title: 'Baseline Latencies',
        path: '/dashboard/cold/baseline',
      },
      {
        title: 'Function Image Size',
        path: '/dashboard/cold/image-size',
      },
      {
        title: 'Deployment Method and Language Runtime',
        path: '/dashboard/cold/deployment-language',
      },
      
    ]
  },
  {
    title: 'Data Transfer Delays',
    path: '/dashboard/data-transfer/',
    icon: getIcon('ci:transfer'),
    children:[
      {
        title:'Inline Transfers',
        path: '/dashboard/data-transfer/inline',
      },
      {
        title:'Storage-based Transfers',
        path: '/dashboard/data-transfer/storage-based',
      },
    ]
  },
  {
    title: 'Bursty Invocations',
    path: '/dashboard/bursty/',
    icon: getIcon('fluent:data-sunburst-24-filled'),
    children:[
      {
        title:'Short IAT',
        path: '/dashboard/bursty/short-iat',
      },
      {
        title:'Long IAT',
        path: '/dashboard/bursty/long-iat',
      },
      {
        title:'Implications of Scheduling Policy',
        path: '/dashboard/bursty/scheduling-policy',
      }
    ]
  },
  {
    title: 'Info',
    path: '/dashboard/info',
    icon: getIcon('eva:file-text-fill'),
  },
];

export default navConfig;
