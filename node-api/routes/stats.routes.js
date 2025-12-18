import { Router } from 'express';
import { calculateStats } from '../controllers/stats.controller.js';

const router = Router();

router.post('/', calculateStats);

export default router;
