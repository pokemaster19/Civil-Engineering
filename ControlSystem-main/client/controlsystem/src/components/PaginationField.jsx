import Pagination from "@mui/material/Pagination";

export const PaginationField = ({count, page, marginBottom, onChange}) => {
    return (
        <Pagination count={count} page={page} onChange={onChange} color="primary" shape="rounded" sx={{ mb: marginBottom }} />
    );
}