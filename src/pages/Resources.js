// material
import { Grid, Card, Container,Link, Stack,Button, Typography, CardContent } from '@mui/material';
// components
import Page from '../components/Page';
// ----------------------------------------------------------------------

export default function Resources() {
  return (
    <Page title="Dashboard: Blog">
      <Container>
        <Stack direction="row" alignItems="center" justifyContent="space-between" mb={5}>
          <Typography variant="h4" gutterBottom>
            Resources
          </Typography>
        </Stack>
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant='h5'>Referencing our work</Typography>
              <br/>
              <Typography variant='subtitle'>If you decide to use STeLLAR for your research and experiments, we are thrilled to support you by offering advice for potential extensions of vHive and always open for collaboration.

<br/><br/>
Please cite  <Link underline="always" color="text.primary" target="_blank" href="/STeLLAR/static/STeLLAR_IISWC21.pdf">
our paper
              </Link> that has recently been accepted to IISWC 2021:
</Typography>
              <Grid container>
              <Grid sx={{opacity:0.6,fontSize:'14px'}} item xs={12} sm={6} md={6}> 
              <br/>{`@inproceedings{ustiugov:analyzing,`}<br/>
{`author    = {Dmitrii Ustiugov and`}<br/>
               {`Theodor Amariucai and`} <br/>
               {`Boris Grot},`} <br/>
  {`title     = {Analyzing Tail Latency in Serverless Clouds with STeLLAR},
  booktitle = {Proceedings of the 2021 IEEE International Symposium on Workload Characterization (IISWC)},`}<br/>
  {`publisher = {{IEEE}},`}
 {`year      = {2021}
}`}
                  </Grid>
              </Grid>
              <Button variant='outlined' target={'_blank'} sx={{my:2}} href='/STeLLAR/static/STeLLAR_IISWC21.pdf'>Download Paper</Button>
              <Button variant='outlined' target={'_blank'} sx={{ml:2,my:2}} href='https://github.com/ease-lab/STeLLAR'>Github</Button>
              
            </CardContent>
          </Card>
        </Grid>
      </Container>
    </Page>
  );
}
