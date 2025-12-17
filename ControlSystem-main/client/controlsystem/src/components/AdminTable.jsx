import { useState, useEffect } from 'react';
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import Button from "@mui/material/Button";
import Box from "@mui/material/Box";
import { EditUserModal } from './Modals';

export const AdminTable = ({ tableWidth, users, page, searchQuery, onUserUpdate, pagination }) => {
    const [openEditDialog, setOpenEditDialog] = useState(false);
    const [currentUser, setCurrentUser] = useState(null);

    const roleMap = {
        1: 'Инженер',
        2: 'Менеджер',
        3: 'Руководитель',
        4: 'Администратор'
    };

    const handleEdit = (user) => {
        setCurrentUser(user);
        setOpenEditDialog(true);
    };

    const handleCloseEditDialog = () => {
        setOpenEditDialog(false);
        setCurrentUser(null);
    };

    return (
        <Box>
            <TableContainer sx={{ width: tableWidth }}>
                <Table stickyHeader>
                    <TableHead>
                        <TableRow>
                            <TableCell>ФИО</TableCell>
                            <TableCell>Почта</TableCell>
                            <TableCell>Роль</TableCell>
                            <TableCell>Статус</TableCell>
                            <TableCell>Действия</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {users.length > 0 ? (
                            users.map((user) => (
                                <TableRow key={user.id}>
                                    <TableCell>{`${user.last_name} ${user.first_name} ${user.middle_name || ''}`}</TableCell>
                                    <TableCell>{user.email}</TableCell>
                                    <TableCell>{roleMap[user.role] || 'Неизвестная роль'}</TableCell>
                                    <TableCell>{user.is_enabled ? 'Активен' : 'Отключён'}</TableCell>
                                    <TableCell>
                                        <Button
                                            variant="contained"
                                            color="primary"
                                            size="small"
                                            onClick={() => handleEdit(user)}
                                        >
                                            Редактировать
                                        </Button>
                                    </TableCell>
                                </TableRow>
                            ))
                        ) : (
                            <TableRow>
                                <TableCell colSpan={5}>Нет пользователей</TableCell>
                            </TableRow>
                        )}
                    </TableBody>
                </Table>
            </TableContainer>

            {currentUser && (
                <EditUserModal
                    open={openEditDialog}
                    user={currentUser}
                    onClose={handleCloseEditDialog}
                    onUserUpdate={onUserUpdate}
                />
            )}
        </Box>
    );
};