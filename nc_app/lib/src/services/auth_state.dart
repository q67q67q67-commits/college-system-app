import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';

class AuthState extends ChangeNotifier {
  String? _token;
  String? _role;
  int? _userId;
  String? _email;
  String? _fullName;

  bool get isAuthenticated => _token != null;
  String? get token => _token;
  String? get role => _role;
  int? get userId => _userId;
  String? get email => _email;
  String? get fullName => _fullName;

  Future<void> loadFromStorage() async {
    final prefs = await SharedPreferences.getInstance();
    final raw = prefs.getString('auth');
    if (raw != null) {
      try {
        final data = jsonDecode(raw) as Map<String, dynamic>;
        _token = data['token'] as String?;
        _role = data['role'] as String?;
        _userId = data['user_id'] is int
            ? data['user_id'] as int
            : (data['user_id'] as num?)?.toInt();
        _email = data['email'] as String?;
        _fullName = data['full_name'] as String?;
      } catch (_) {
        await clear();
      }
    }
    notifyListeners();
  }

  Future<void> saveAuth(Map<String, dynamic> payload) async {
    _token = payload['token'] as String?;
    _role = payload['role'] as String?;
    _userId = payload['user_id'] is int
        ? payload['user_id'] as int
        : (payload['user_id'] as num?)?.toInt();
    _email = payload['email'] as String?;
    _fullName = payload['full_name'] as String?;
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('auth', jsonEncode({
      'token': _token,
      'role': _role,
      'user_id': _userId,
      'email': _email,
      'full_name': _fullName,
    }));
    notifyListeners();
  }

  Future<void> clear() async {
    _token = null;
    _role = null;
    _userId = null;
    _email = null;
    _fullName = null;
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('auth');
    notifyListeners();
  }
}
