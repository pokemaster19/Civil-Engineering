import Typography from "@mui/material/Typography";
import { Grid, Box } from "@mui/material";

export const DefectCounter = ({total, inProgress, opened, resolved, overdued}) => {
    return (
        <Grid container spacing={3}>
            <Box sx={{display: {xs: 'none', sm: 'block', xl: 'block', lg: 'block'}}}>
                <Typography>Всего</Typography>
                <Typography variant="h4">{total}</Typography>
            </Box>
            <Box>
                <Typography color="#00BBFF">Открытые</Typography>
                <Typography variant="h4" color="#00BBFF">{opened}</Typography>
            </Box>
            <Box>
                <Typography color="#F8BE00">В работе</Typography>
                <Typography variant="h4" color="#F8BE00">{inProgress}</Typography>
            </Box>
            <Box>
                <Typography color="#6FDF00">Исправлены</Typography>
                <Typography variant="h4" color="#6FDF00">{resolved}</Typography>
            </Box>
            <Box sx={{display: {xs: 'none', sm: 'block', xl: 'block', lg: 'block'}}}>
                <Typography color="red">Просрочены</Typography>
                <Typography variant="h4" color="red">{overdued}</Typography>
            </Box>
        </Grid>
    );
}
