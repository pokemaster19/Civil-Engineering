import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import {AuthProvider} from './context/AuthContext';
import LoginPage from './pages/LoginPage';
import ProjectsPage from './pages/ProjectsPage';
import AdminPage from './pages/AdminPage';
import ProjectPage from './pages/ProjectPage';
import DefectPage from './pages/DefectPage';
import ProtectedRoute from './components/ProtectedRoute';
import { GuestRoute } from './context/GuestContext';
import { AdminRoute } from './components/AdminRoute';

function App() {
  return (
    <Router>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<GuestRoute><LoginPage /></GuestRoute>} />
          <Route path='/' element={<ProtectedRoute><ProjectsPage /></ProtectedRoute>} />
          <Route path='/project/:projectId' element={<ProtectedRoute><ProjectPage /></ProtectedRoute>} />
          <Route path='/admin' element={<AdminRoute><AdminPage /></AdminRoute>} />
          <Route path='/defect/:defectId' element={<ProtectedRoute><DefectPage /></ProtectedRoute>} />
        </Routes>
      </AuthProvider>
    </Router>
  );
}

export default App;
