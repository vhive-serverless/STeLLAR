// component
import Iconify from '../../components/Iconify';

// ----------------------------------------------------------------------

const getIcon = (name) => <Iconify icon={name} width={22} height={22} />;

const navConfig = [
  {
    title: 'About',
    path: '/dashboard/about',
    icon: getIcon('mdi:graph-line'),
  },
  {
    title: 'Warm Function Invocations',
    path: '/dashboard/warm/aws',
    icon: getIcon('bxs:hot'),
  },
  {
    title: 'Cold Function Invocations',
    path: '/dashboard/cold/baseline',
    icon: getIcon('ic:baseline-severe-cold'),
  },
  {
    title: 'Cold Function - Image Size',
    path: '/dashboard/cold/image-size',
    icon: getIcon('ic:baseline-severe-cold'),
  },
  {
    title: 'Cold Function - Deployment and Language Runtime',
    path: '/dashboard/cold/deployment-language',
    icon: getIcon('ic:baseline-severe-cold'),
  },
  // {
  //   title: 'Data Transfer Delays',
  //   path: '/dashboard/data-transfer/',
  //   icon: getIcon('ci:transfer'),
  //   children:[
  //     {
  //       title:'Inline Transfers',
  //       path: '/dashboard/data-transfer/inline',
  //     },
  //     {
  //       title:'Storage-based Transfers',
  //       path: '/dashboard/data-transfer/storage-based',
  //     },
  //   ]
  // },
  // {
  //   title: 'Bursty Invocations',
  //   path: '/dashboard/bursty/',
  //   icon: getIcon('fluent:data-sunburst-24-filled'),
  //   children:[
  //     {
  //       title:'Short IAT',
  //       path: '/dashboard/bursty/short-iat',
  //     },
  //     {
  //       title:'Long IAT',
  //       path: '/dashboard/bursty/long-iat',
  //     },
  //     {
  //       title:'Implications of Scheduling Policy',
  //       path: '/dashboard/bursty/scheduling-policy',
  //     }
  //   ]
  // },
  {
    title: 'Team',
    path: '/dashboard/team',
    icon: getIcon('ri:team-fill'),
  },
];

export default navConfig;
