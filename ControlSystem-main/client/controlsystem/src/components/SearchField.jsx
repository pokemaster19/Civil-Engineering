import SearchIcon from '@mui/icons-material/Search';
import TextField from "@mui/material/TextField";
import IconButton from '@mui/material/IconButton';
import InputAdornment from '@mui/material/InputAdornment';

export const SearchField = ({ value, onChange, onSearchClick }) => {
    return (
        <TextField
            id="outlined-search"
            label="Найти"
            type="search"
            value={value}
            onChange={(e) => onChange(e.target.value)}
            InputProps={{
                endAdornment: (
                    <InputAdornment position="end">
                        <IconButton edge="end" aria-label="search" onClick={onSearchClick}>
                            <SearchIcon />
                        </IconButton>
                    </InputAdornment>
                ),
            }}
            sx={{ width: { md: '50vw', xs: '85vw' } }}
        />
    );
};
