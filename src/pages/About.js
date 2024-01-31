// material
import { Grid, Card, Container, Stack,Box,Button, Typography, CardContent,Link } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import ListItem from '@mui/material/ListItem';
import {ReactComponent as ArchitectureDiagram} from '../images/diagram.svg';
// components
import Page from '../components/Page';

export default function About() {
  return (
    <Page title="Dashboard: About">
      <Container maxWidth="xl">
        <Stack direction="row" alignItems="center" justifyContent="space-between" mt={3}>
          {/* <Typography variant="h4" gutterBottom>
          About
          </Typography> */}
        </Stack>

        <Grid item xs={12} mb={3}>
          <Card>
            <CardContent>
            <Typography variant='h4' marginBottom={3}>What is STeLLAR ? (Serverless Tail-Latency Analyzer)</Typography>
              <Typography variant='p'>STeLLAR is an open-source serverless benchmarking framework, which enables an accurate performance characterization of serverless deployments. 
              STeLLAR is provider-agnostic and highly configurable, allowing the analysis of both end-to-end and per-component performance with minimal instrumentation effort. 
              Using STeLLAR, we continously conduct various performance tests in different serverless settings and these results are visualized in a dashboard. <br/>
              STeLLAR is a part of the <Link target="_blank" href={'https://vhive-serverless.github.io/'}>vHive Ecosystem.</Link>
 
</Typography>
            </CardContent>



          </Card>
          
        </Grid>

        <Box item xs={12} mb={3}>
          <Card>
            <CardContent>
            <Typography variant='h4' marginBottom={3}>Design & Methodology</Typography>
            {/* <Stack direction="row" alignItems="center" mt={5}> */}
            <Box sx={{display: 'flex',my:3,alignItems:'center',justifyContent:'center'}}>
            
            <ArchitectureDiagram width={'90%'}/>
            </Box>
            <Box sx={{display: 'flex',my:3,alignItems:'center',justifyContent:'center' }}>
            STeLLAR Architecture Overview
              </Box>
            <Box sx={{ width: '100%',ml:5}}>
              
            <Typography variant='h6'><b>Terminology</b></Typography>
            
            <ListItem sx={{display:'list-item'}}>
       An endpoint is a URL used for locating the function instance over the Internet. As seen in the diagram, this URL most often points to resources such as AWS API Gateway, Azure HTTP Triggers, or similar.
       </ListItem>

       <ListItem sx={{display:'list-item'}}>The inter-arrival time (IAT) is the time interval that the client waits for in-between sending two bursts to the same endpoint.
</ListItem>


<Typography variant='h6' mt={2}><b>Components</b></Typography>
          <ListItem sx={{ display: 'list-item' }}>
          The client consists of two primary elements: the deployer and the load generator. 
          The deployer is responsible for coordinating the deployment of functions to various cloud providers, while the load generator is tasked with sending requests to the deployed functions.


          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
          The experiment configuration is an input JSON file used to specify and customize the experiments.

          </ListItem>
      

       
       <ListItem sx={{display:'list-item'}}>The latencies CSV files are the main output of the evaluation framework. They are used in our plotting tools to produce insightful visualizations.
</ListItem>

<ListItem sx={{display:'list-item'}}>The logs text file is the final output of the benchmarking client. Log records are useful for optimizing code and debugging problematic behavior.</ListItem>

              </Box>
              
            {/* </Stack> */}
            <Typography variant='p' color={'red'}>* Currently we only supports visualizing benchmarking results from AWS Lambda, Google Cloud Run, Azure and Cloudflare but we will extend to other cloud providers in the future.</Typography>
            </CardContent>
            <CardContent>
            <Typography variant='h5'>Client Configuration</Typography>
            <br/>
            <Typography variant={'p'} sx={{ mb: 2 }}>
               <b>AWS</b>
            </Typography>

      <Stack direction="row" alignItems="center" mt={2}>
            <Box sx={{ width: '100%',ml:1}}>
          <ListItem sx={{ display: 'list-item' }}>
            Instance Type : <b>t2.micro</b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            CPU : <b>1 vCPUs</b>
          </ListItem>
         
         

          </Box>
            <Box sx={{ width: '100%',ml:1}}>
            <ListItem sx={{ display: 'list-item' }}>
            Datacenter : <b>us-west-1</b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            Memory : <b>1.0 GiB RAM</b>
          </ListItem>
         
              </Box>
              
              </Stack>
            
              <br/>
            <Typography variant={'p'} sx={{ mb: 2 }}>
               <b>Azure</b>
            </Typography>

<Stack direction="row" alignItems="center" mt={2}>
            <Box sx={{ width: '100%',ml:1}}>
          <ListItem sx={{ display: 'list-item' }}>
            Instance Type : <b>Standard B1s</b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            CPU : <b>1 vCPUs</b>
          </ListItem>
         
         

          </Box>
            <Box sx={{ width: '100%',ml:1}}>
            <ListItem sx={{ display: 'list-item' }}>
            Datacenter : <b>West US</b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            Memory : <b>1.0 GiB RAM</b>
          </ListItem>
         
              </Box>
              </Stack>
            

              <br/>
            <Typography variant={'p'} sx={{ mb: 2 }}>
               <b>Google Cloud</b>
            </Typography>

<Stack direction="row" alignItems="center" mt={2}>
            <Box sx={{ width: '100%',ml:1}}>
          <ListItem sx={{ display: 'list-item' }}>
            Instance Type : <b>Compute Engine e2-micro</b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            CPU : <b>2 vCPUs</b>
          </ListItem>
         
         

          </Box>
            <Box sx={{ width: '100%',ml:1}}>
            <ListItem sx={{ display: 'list-item' }}>
            Datacenter : <b>us-west1-a</b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            Memory : <b>1.0 GiB RAM</b>
          </ListItem>
         
              </Box>
              </Stack>

              <br/>
            <Typography variant={'p'} sx={{ mb: 2 }}>
               <b>Cloudflare</b>
            </Typography>

<Stack direction="row" alignItems="center" mt={2}>
            <Box sx={{ width: '100%',ml:1}}>
          <ListItem sx={{ display: 'list-item' }}>
            Instance Type : <b>AWS t2.micro</b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            CPU : <b>1 vCPUs</b>
          </ListItem>
         
         

          </Box>
            <Box sx={{ width: '100%',ml:1}}>
            <ListItem sx={{ display: 'list-item' }}>
            Datacenter : <b>us-east-2</b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            Memory : <b>1.0 GiB RAM</b>
          </ListItem>
         
              </Box>
              </Stack>

             <br/>
               <ListItem sx={{display:'list-item'}}>
               <b>Scheduled Experiment time: 00:00 UTC</b>
            </ListItem>
            {/* <ListItem sx={{display:'list-item'}}>We collect <b>1000 samples</b> for all experiments except for warm function experiment where we collect <b>3000 samples</b> to generate statistics.
            </ListItem> */}

            </CardContent>
          
            <CardContent>
            <Typography variant='h5' marginBottom={2}>Function Deployment Configuration</Typography>
              <ListItem sx={{display:'list-item'}}>
            Our experiments are based on <b>Python 3 functions</b>, except for evaluations of language runtimes where we evaluate four different language runtimes. </ListItem>
            <ListItem sx={{display:'list-item'}}>We commonly employ <b>ZIP-based deployment</b> as our default method of deployment, with the exception of Google Cloud Run (GCR), which exclusively requires <b>Container based deployment.</b></ListItem><br/>
       <b>Function Deployment Regions</b>
              <ListItem sx={{display:'list-item'}}>AWS - us-west-1</ListItem>
              <ListItem sx={{display:'list-item'}}>Azure - West US</ListItem>
              <ListItem sx={{display:'list-item'}}>Google Cloud - us-west1-a</ListItem>
              <ListItem sx={{display:'list-item'}}>Cloudflare - n/a</ListItem>
            </CardContent>
          </Card>
          
        </Box>

        <Grid item xs={12} mb={3}>
          <Card>
            <CardContent>
            <Typography variant='h4'>Scenarios under test</Typography>
              <Typography variant='p'>There are several important scenarios under test and displayed in the dashboards, including the following:</Typography>
            </CardContent>

            <CardContent>
            <Typography variant='h5' fontWeight={500}>1. Warm Function Invocations <Button to="/dashboard/warm/aws" size="small" variant="outlined" sx={{marginLeft:3,color:'green'}} component={RouterLink}>
            View Results
          </Button></Typography>
            <ListItem sx={{display:'list-item'}}>Under warm function invocations, we evaluate the response time of warm functions under individual invocations (i.e., allowing no more than a single outstanding request to each function).
<br/>
            </ListItem>
            
              <Typography variant='h5' fontWeight={500} mt={3}>2. Cold Function Invocations</Typography>
              Under Cold Function Invocation experiments, we evaluate the response time of functions with cold instances by issuing invocations (one-request at a time) with a long inter-arrival time (IAT) of 600 seconds.
              <ListItem sx={{display:'list-item'}}><b>Basic :</b> We evaluate the response time of functions with cold instances.  
              <Button to="/dashboard/cold/baseline" size="small" variant="outlined" sx={{marginLeft:3,color:'green'}} component={RouterLink}>
            View Results
          </Button>
          </ListItem>
              <ListItem sx={{display:'list-item'}}><b>Function Container Image Size : </b> We evaluate the impact of container image size on response time for functions with cold instances.
              <Button to="/dashboard/cold/image-size" size="small" variant="outlined" sx={{marginLeft:3,color:'green'}} component={RouterLink}>
            View Results
          </Button>
          
          {/* <ListItem sx={{ml:3,opacity: 0.7, display:'list-item'}}>Image Sizes : 50MB , 100MB</ListItem> */}
       
          </ListItem>
          <ListItem sx={{display:'list-item'}}><b>Language Runtime : </b> We evaluate the impact of different language runtimes on response time for functions with cold instances.
              <Button to="/dashboard/cold/deployment-language" size="small" variant="outlined" sx={{marginLeft:3,color:'green'}} component={RouterLink}>
                  View Results
              </Button>
          </ListItem>
          <Typography variant='p'>Disclaimer: Experiments marked 'No data' reflect repeated unsuccessful attempts.</Typography>
          {/* <Typography variant='p'>* We call a function warm if it has at least one instance online and idle upon a request’s arrival, otherwise we refer to the function as a cold function.</Typography> */}
          </CardContent>
          </Card>
          
        </Grid>
      </Container>
    </Page>
  );
}
