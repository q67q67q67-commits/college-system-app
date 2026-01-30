import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../services/api_service.dart';

class AuthProvider with ChangeNotifier {
  String? _token;
  int? _userId;
  String? _role;
  String? _email;
  String? _fullName;

  String? get token => _token;
  int? get userId => _userId;
  String? get role => _role;
  String? get email => _email;
  String? get fullName => _fullName;
  bool get isLoggedIn => _token != null && _token!.isNotEmpty;

  static const _keyToken = 'nc_token';
  static const _keyUserId = 'nc_user_id';
  static const _keyRole = 'nc_role';
  static const _keyEmail = 'nc_email';
  static const _keyFullName = 'nc_full_name';

  Future<void> loadFromStorage() async {
    final prefs = await SharedPreferences.getInstance();
    _token = prefs.getString(_keyToken);
    _userId = prefs.getInt(_keyUserId);
    _role = prefs.getString(_keyRole);
    _email = prefs.getString(_keyEmail);
    _fullName = prefs.getString(_keyFullName);
    if (_token != null) ApiService.setToken(_token);
    notifyListeners();
  }

  Future<void> login(String email, String password) async {
    final res = await ApiService.post('/auth/login', {'email': email, 'password': password});
    _token = res['token'] as String?;
    _userId = (res['user_id'] as num?)?.toInt();
    _role = res['role'] as String?;
    _email = res['email'] as String?;
    _fullName = res['full_name'] as String?;
    ApiService.setToken(_token);
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_keyToken, _token ?? '');
    await prefs.setInt(_keyUserId, _userId ?? 0);
    await prefs.setString(_keyRole, _role ?? '');
    await prefs.setString(_keyEmail, _email ?? '');
    await prefs.setString(_keyFullName, _fullName ?? '');
    notifyListeners();
  }

  Future<void> logout() async {
    _token = _userId = null;
    _role = _email = _fullName = null;
    ApiService.setToken(null);
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_keyToken);
    await prefs.remove(_keyUserId);
    await prefs.remove(_keyRole);
    await prefs.remove(_keyEmail);
    await prefs.remove(_keyFullName);
    notifyListeners();
  }

  bool get isTeacher => _role == 'teacher' || _role == 'director' || _role == 'admin';
  bool get isAdminOrDirector => _role == 'admin' || _role == 'director';
  bool get showDirectorProfile => _role != 'admin' && _role != null;
}
