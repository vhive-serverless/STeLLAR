// material
import { Grid, Card,Box, Container, Stack,Typography,Link } from '@mui/material';
// components
import Page from '../components/Page';
import AccountProfile from '../components/AccountProfile'
// data
import { members,alumni } from '../data/members';
// ----------------------------------------------------------------------

const styles = {
  container: {
    display: 'flex',
    justifyContent: 'center',
    alignItems:'center'
  
  },
  card: {
    width:'100%',
  },
};



export default function Info() {


  return (
    <Page title="Dashboard: Info" sx={styles.container}>
      <Container maxWidth="xl">
        
        <Stack direction="row" alignItems="center" justifyContent="space-between" mt={3}>
          <Typography variant="h4" gutterBottom>
            Team
          </Typography>
        </Stack>
          <Grid container xs={12}>
                {members.map((member) => (
                   <Grid item lg={3}
                   md={6}
                   xs={12}
                   sx={{pr:2,pb:2}} key={member.name} >
                    <Card>
            <AccountProfile member={member}/>
            </Card>
            </Grid>
            ))}

            </Grid>

            <Stack direction="row" alignItems="center" justifyContent="space-between" mt={3}>
          <Typography variant="h4" gutterBottom>
            Alumni
          </Typography>
        </Stack>
          <Grid container xs={12}>
                {alumni.map((member) => (
                   <Grid item lg={3}
                   md={6}
                   xs={12}
                   sx={{pr:2,pb:2}} key={member.name} >
                    <Card>
            <AccountProfile member={member}/>
            </Card>
            </Grid>
            ))}

            </Grid>
        <Card sx={styles.card} > 

          <Stack direction="row" alignItems="center" justifyContent="center" sx={{paddingY:0.9}}>                 
          <Typography
            color="textPrimary"
            sx={{fontSize:16,fontWeight:500,marginX:2}}
            align='center'
          >
             Find us on :
          </Typography>     
            <Link target="_blank" href={'https://github.com/vhive-serverless'}><Box component="img" src="/STeLLAR/static/icons/github.png" sx={{height: '30px',marginX:2 }}/></Link>
            <Link target="_blank" href={'https://www.youtube.com/playlist?list=PLVdxPJaekjWqBsEUwnrYRQCaMqvcDVsBE'}><Box component="img" src="/STeLLAR/static/icons/youtube.png" sx={{height: '30px',marginX:2  }}/></Link>
            <Link target="_blank" href={'https://join.slack.com/t/vhivetutorials/shared_invite/zt-1fk4v71gn-nV5oev5sc9F4fePg3_OZMQ'}><Box component="img" src="/STeLLAR/static/icons/slack.png" sx={{height: '30px',marginX:2  }}/></Link>
          </Stack>

            </Card>
      </Container>
    </Page>
  
  );
}
