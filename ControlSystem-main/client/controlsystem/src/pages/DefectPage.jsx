import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import Typography from '@mui/material/Typography';
import {Header} from '../components/AppBar';
import background from "../css/Background.module.css";
import { ReportsTable } from '../components/ReportsTable';
import { PreviewTable } from '../components/PrevievTable';
import { Box, CircularProgress, Alert, Button } from '@mui/material';
import { fetchDefectById, leaveComment, fetchComments } from '../api/Defects';
import { EditEntityModal } from '../components/Modals';
import { uploadAttachment } from '../api/Attachments';

const DefectPage = () => {
    const { defectId } = useParams();
    const [defect, setDefect] = useState(null);
    const [comments, setComments] = useState([]);
    const [files, setFiles] = useState([]);
    const [pagination, setPagination] = useState({
        page: 1,
        limit: 10,
        total: 0,
        totalPages: 0,
        hasNextPage: false,
        hasPrevPage: false,
    });
    const [loading, setLoading] = useState(true);
    const [commentsLoading, setCommentsLoading] = useState(false);
    const [error, setError] = useState(null);

    const loadComments = async (page = 1, limit = 10) => {
        try {
            setCommentsLoading(true);
            setError(null);
            const commentsResponse = await fetchComments(defectId, page, limit);
            setComments(commentsResponse.comments || []);
            setPagination(commentsResponse.pagination || {
                page: 1,
                limit: 10,
                total: 0,
                totalPages: 0,
                hasNextPage: false,
                hasPrevPage: false,
            });
        } catch (err) {
            setError('Failed to load comments: ' + err.message);
        } finally {
            setCommentsLoading(false);
        }
    };

    const loadInitialData = async () => {
        try {
            setLoading(true);
            setError(null);
            const defectResponse = await fetchDefectById(defectId);
            setDefect(defectResponse.defect);
            await loadComments(1, pagination.limit);
        } catch (err) {
            setError('Failed to load data: ' + err.message);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (defectId) {
            loadInitialData();
        }
    }, [defectId]);

    const handleCommentSubmit = async (comment) => {
        try {
            const newComment = await leaveComment(defectId, comment);
            setComments((prevComments) => [newComment, ...prevComments]);
            setPagination((prev) => ({
                ...prev,
                total: prev.total + 1,
                totalPages: Math.ceil((prev.total + 1) / prev.limit),
                hasNextPage: prev.page < Math.ceil((prev.total + 1) / prev.limit),
            }));
        } catch (err) {
            setError('Failed to post comment: ' + err.message);
        }
    };

    const handlePageChange = (newPage) => {
        loadComments(newPage, pagination.limit);
    };

    const handleFileChange = (event) => {
        setFiles(Array.from(event.target.files));
    };

    const handleUploadFiles = async () => {
        if (files.length === 0) {
            alert('Пожалуйста, выберите хотя бы один файл.');
            return;
        }

        try {
            for (const file of files) {
                const formData = new FormData();
                formData.append('file', file);
                formData.append('defectId', defectId);
                await uploadAttachment(formData);
            }
            setFiles([]);
            setError(null);
            const defectResponse = await fetchDefectById(defectId);
            setDefect(defectResponse.defect);
            alert('Вложения успешно загружены!');
        } catch (err) {
            setError('Ошибка при загрузке вложений: ' + err.message);
        }
    };


    if (loading) {
        return (
        <div className={background.background}>
            <Header />
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
            <CircularProgress />
            </Box>
        </div>
        );
    }

    if (error) {
        return (
        <div className={background.background}>
            <Header />
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
            <Alert severity="error">{error}</Alert>
            </Box>
        </div>
        );
    }

    if (!defect) {
        return (
        <div className={background.background}>
            <Header />
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
            <Alert severity="warning">Defect not found</Alert>
            </Box>
        </div>
        );
    }
    return (
        <div className={background.background}>
            <Header />
            <div className={background.contentParent}>
                <Box sx={{
                display: 'flex',
                flexDirection: 'column',
                gap: '30px',
                }}>
                    <Typography variant='h3'>{defect.title}</Typography>
                    <Typography variant='body1'>{defect.description}</Typography>
                    <EditEntityModal entityType={"defect"} entityId={defectId} title={defect.title} description={defect.description} status={defect.status} priority={defect.priority}></EditEntityModal>
                    <Typography variant='h4'>Медиафайлы</Typography>
                    <PreviewTable images={defect.photosUrl || []} files={defect.filesUrl || []} />
                    <Box sx={{ mb: 2 }}>
                        <Typography variant="subtitle1" sx={{ mb: 1 }}>Загрузка новых вложений</Typography>
                        <input
                            type="file"
                            accept="image/*,application/pdf"
                            onChange={handleFileChange}
                            style={{ marginBottom: '8px' }}
                            multiple
                        />
                        <Button variant="contained" onClick={handleUploadFiles}>
                            Загрузить файлы
                        </Button>
                    </Box>
                    <Typography variant='h4'>Комментарии</Typography>
                    {commentsLoading ? (
                        <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '200px' }}>
                            <CircularProgress />
                        </Box>
                    ) : (
                        <ReportsTable
                            onCommentSubmit={handleCommentSubmit}
                            comments={comments}
                            pagination={pagination}
                            onPageChange={handlePageChange}
                        />
                    )}
                </Box>
            </div>
        </div>
    );
}

export default DefectPage