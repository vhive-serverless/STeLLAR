// material
import { Grid, Card, Container, Stack,Box,Button, Typography, CardContent } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import ListItem from '@mui/material/ListItem';
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

            <Typography variant='h4' marginBottom={3}>Serverless Computing & Benchmarking</Typography>
              <Typography variant='p'>Serverless computing, also known as Function-as-a-
Service (FaaS), where developers structure their applications as a collection of functions invoked by various
events like clicks, and cloud providers take responsibility for
cloud infrastructure management. As with other cloud services,
serverless deployments require responsiveness and performance
predictability manifested through low average and tail latencies. Hence, we introduce STeLLAR : A framework to
address the lack of a toolchain for tail-latency analysis in serverless computing.
 
</Typography>
            <Typography variant='h4' marginTop={5} marginBottom={3}>What is STeLLAR ? (Serverless Tail-Latency Analyzer)</Typography>
              <Typography variant='p'>STeLLAR is an open-source serverless benchmarking framework, which enables an accurate performance characterization of serverless deployments. 
              STeLLAR is provider-agnostic and highly configurable, allowing the analysis of both end-to-end and per-component performance with minimal instrumentation effort. 
              Using STeLLAR, we continously conduct various performance tests in different serverless settings and these results are visualized in a dashboard.
 
</Typography>
            </CardContent>



          </Card>
          
        </Grid>

        <Box item xs={12} mb={3}>
          <Card>
            <CardContent>
            <Typography variant='h4' marginBottom={3}>Design & Methodology</Typography>
            {/* <Stack direction="row" alignItems="center" mt={5}> */}
            <Box sx={{display: 'flex',my:3,alignItems:'center',justifyContent:'center' }}>
            <Box component="img" src="/STeLLAR/static/design.png" sx={{width:'60%'}} />
            </Box>
            <Box sx={{display: 'flex',my:3,alignItems:'center',justifyContent:'center' }}>
            Figure 1: STeLLAR Architecture Overview
              </Box>
            <Box sx={{ width: '100%',ml:5}}>
              
            <Typography variant='h6'><b>Terminology</b></Typography>
            
            <ListItem sx={{display:'list-item'}}>
       An endpoint is a URL used for locating the function instance over the Internet. As seen in the diagram, this URL most often points to resources such as AWS API Gateway, Azure HTTP Triggers, vHive Kubernetes Load Balancer, or similar.
       </ListItem>

       <ListItem sx={{display:'list-item'}}>The inter-arrival time (IAT) is the time interval that the client waits for in-between sending two bursts to the same endpoint.
</ListItem>


       <ListItem sx={{display:'list-item'}}>Multiple endpoints can be used simultaneously by the same experiment to speed up the benchmarking. The JSON configuration field parallelism defines this number: the higher it is, the more endpoints will be allocated, and the more bursts will be sent in short succession (speeding up the process for large IATs).
</ListItem>

<Typography variant='h6' mt={2}><b>Components</b></Typography>
          <ListItem sx={{ display: 'list-item' }}>
          The coordinator orchestrates the entire benchmarking procedure.


          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
          The experiment configuration is an input JSON file used to specify and customize the experiments.

          </ListItem>
       
      
              
       <ListItem sx={{display:'list-item'}}>The vendor endpoints input JSON file is only used for providers such as vHive that do not currently support automated function management (e.g., function listing, deployment, repurposing, or removal via SDKs or APIs).
</ListItem>

       
       <ListItem sx={{display:'list-item'}}>The latencies CSV files are the main output of the evaluation framework. They are used in our plotting tools to produce insightful visualizations.
</ListItem>

<ListItem sx={{display:'list-item'}}>The logs text file is the final output of the benchmarking client. Log records are useful for optimizing code and debugging problematic behavior.</ListItem>

              </Box>
              
            {/* </Stack> */}
            <Typography variant='p' color={'red'}>* Currently we only supports visualizing benchmarking results from AWS Lambda</Typography>
            </CardContent>
            <CardContent>
            <Typography variant='h5' marginBottom={2}>Client Configuration</Typography>
            <ListItem sx={{display:'list-item'}}>We run the STeLLAR client on <b>t2.small</b> node in <b>AWS - Oregon (us-west-2)</b> datacenter region which features a <b>Intel Xeon CPU</b> with <b>2GB DRAM</b>.</ListItem>
            <ListItem sx={{display:'list-item'}}>We initiate the experiments sequentially at <b> 00:00h (GMT) </b> on each day.
            </ListItem>
            
            <ListItem sx={{display:'list-item'}}>We collect <b>1000 samples</b> for all experiments except for warm function experiment where we collect <b>3000 samples</b> to generate statistics.
            </ListItem>

            </CardContent>
          
            <CardContent>
            <Typography variant='h5' marginBottom={2}>Function Deployment Configuration</Typography>
              <ListItem sx={{display:'list-item'}}>We deploy the functions in same datacenter region where STeLLAR client runs which is AWS - Oregon (us-west-2). </ListItem>
              <ListItem sx={{display:'list-item'}}>The functions are configured with the different memory sizes, ranging from <b>128MB to 2GB </b> and are specified in the section of each experiment.</ListItem>
              <ListItem sx={{display:'list-item'}}>Unless specified otherwise, we deploy all functions using the <b>ZIP-based deployment</b> method and use <b>Python 3</b> functions for all experiments except the function image size experiment. <br/>In this experiments, we use <b>Golang</b> functions to minimize the image size.
 
            </ListItem>
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
            <ListItem sx={{display:'list-item'}}>Under warm function invocations, we evaluate the response time of warm functions * under a non-bursty load (i.e., allowing no more than a single outstanding request to each function).
<br/>
            </ListItem>
            
              <Typography variant='h5' fontWeight={500} mt={3}>2. Cold Function Invocations</Typography>
              <ListItem sx={{display:'list-item'}}><b>Basic :</b> Here we assess the cold function delays under a non-bursty load.  
              <Button to="/dashboard/cold/baseline" size="small" variant="outlined" sx={{marginLeft:3,color:'green'}} component={RouterLink}>
            View Results
          </Button>
          </ListItem>
              <ListItem sx={{display:'list-item'}}><b>Function Image Size : </b> Secondly, we evaluate response times of cold functions varying the function image sizes in 3 different settings.
              <Button to="/dashboard/cold/image-size" size="small" variant="outlined" sx={{marginLeft:3,color:'green'}} component={RouterLink}>
            View Results
          </Button>
          
          <ListItem sx={{ml:3,opacity: 0.7, display:'list-item'}}>Image Size s : 10MB , 60MB , 100MB</ListItem>
       
          </ListItem>
          <ListItem sx={{display:'list-item'}}><b>Deployment Method & Language Runtime : </b> In here, we measure the response times of cold functions varying the language runtime and the deployment method. 
              <Button to="/dashboard/cold/deployment-language" size="small" variant="outlined" sx={{marginLeft:3,color:'green'}} component={RouterLink}>
                  View Results
              </Button>
            <ListItem sx={{ml:3,opacity: 0.7, display:'list-item'}}>Python - ZIP based deployment</ListItem>
            <ListItem sx={{ml:3,opacity: 0.7, display:'list-item'}}>Python - Image based deployment</ListItem>
            <ListItem sx={{ml:3,opacity: 0.7, display:'list-item'}}>Go - ZIP based deployment</ListItem>
            <ListItem sx={{ml:3,opacity: 0.7, display:'list-item'}}>Go - Image based deployment</ListItem>
          </ListItem>
          <Typography variant='p'>* We call a function warm if it has at least one instance online and idle upon a request’s arrival, otherwise we refer to the function as a cold function.</Typography>
          </CardContent>
          </Card>
          
        </Grid>
      </Container>
    </Page>
  );
}
