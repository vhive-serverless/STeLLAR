import PropTypes from 'prop-types';
import {
    Avatar,
    Box,
    Container,
    Divider,
    Link,
    Typography
  } from '@mui/material';

 
  const AccountProfile = ({member}) => (
    <Container>
     
        <Box
          sx={{
            alignItems: 'center',
            justifyContent:'center',
            display: 'flex',
            flexDirection: 'column',
            padding:2
          }}
        >
          
          <Avatar
            src={member.photo}
            sx={{
              height: 64,
              mb: 2,
              width: 64,
            }}
          />
          <Link target="_blank" href={member.link}>
          <Typography
            color="textPrimary"
            gutterBottom
            sx={{fontSize:14,fontWeight:600}}
            align='center'
          >
            {member.name}
          </Typography>
          </Link>
          <Typography
            color="textSecondary"
            justifyContent={'center'}
            sx={{fontSize:12,textAlign:'center'}}

          >
            {member.affiliation}
          </Typography>

        </Box>
      <Divider />
    </Container>
  );
  
  AccountProfile.propTypes = {
    member: PropTypes.object,
  };
  export default AccountProfile;