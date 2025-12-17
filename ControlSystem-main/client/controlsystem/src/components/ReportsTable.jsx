import Box from "@mui/material/Box";
import { useState } from "react";
import Input from "@mui/material/Input";
import InputAdornment from "@mui/material/InputAdornment";
import { PaginationField } from "./PaginationField";
import SendIcon from '@mui/icons-material/Send';
import Typography from "@mui/material/Typography";


const ReportCard = ({authorName, content}) => {
    return(
        <Box sx={{
            background: 'white',
            borderRadius: '5px',
            padding: '5px 5px 5px 5px',
        }}>
            <Typography variant="h6">{authorName}</Typography>
            <Typography variant="body2" color="grey">{content}</Typography>
        </Box>
    );
}

export const ReportsTable = ({ onCommentSubmit, comments = [], pagination, onPageChange }) => {
  const [comment, setComment] = useState('');

  const handleSubmit = () => {
    if (comment.trim()) {
      onCommentSubmit(comment);
      setComment('');
    }
  };

  return (
    <Box sx={{
      minWidth: '80vw',
      display: 'flex',
      flexDirection: 'column',
      gap: '15px',
      background: '#f3f3f3ff',
      padding: '20px 20px 20px 20px',
    }}>
      <Input
        id="reports-input"
        multiline
        maxRows={5}
        placeholder="Комментарий"
        fullWidth
        value={comment}
        onChange={(e) => setComment(e.target.value)}
        endAdornment={
          <InputAdornment position="end">
            <SendIcon onClick={handleSubmit} style={{ cursor: 'pointer' }} />
          </InputAdornment>
        }
      />
      {comments.map((comment) => (
        <ReportCard
          key={comment.id}
          authorName={comment.authorName}
          content={comment.content}
        />
      ))}
      <PaginationField 
        page={pagination.page}
        count={pagination.totalPages} 
        onChange={(e, value) => onPageChange(value)}
      />
    </Box>
  );
};