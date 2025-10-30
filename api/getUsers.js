const express = require('express');
const admin = require('firebase-admin');

const router = express.Router();

// Firebase Admin SDK 초기화
admin.initializeApp({
  credential: admin.credential.cert(require('../electric-circuit-web-firebase-adminsdk-fbsvc-daa68ebecf.json')),
});

// 사용자 목록 가져오기 API
router.get('/users', async (req, res) => {
  try {
    const result = await admin.auth().listUsers(1000);
    const users = result.users.map((userRecord) => userRecord.toJSON());
    res.status(200).json(users);
  } catch (error) {
    console.error('Error fetching users:', error);
    res.status(500).json({ error: 'Failed to fetch users' });
  }
});

module.exports = router;