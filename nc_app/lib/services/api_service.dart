import 'dart:convert';
import 'package:http/http.dart' as http;
import '../config/api_config.dart';

class ApiService {
  static String? _token;
  static String get baseUrl => apiBaseUrl;
  /// Вызывается при 401 — приложение должно выйти и открыть экран входа.
  static void Function()? onUnauthorized;

  static void setToken(String? t) => _token = t;

  static Map<String, String> get _headers {
    final m = {'Content-Type': 'application/json'};
    if (_token != null) m['Authorization'] = 'Bearer $_token';
    return m;
  }

  static Future<Map<String, dynamic>> get(String path) async {
    final r = await http.get(Uri.parse('$baseUrl$path'), headers: _headers);
    return _handle(r);
  }

  static Future<Map<String, dynamic>> post(String path, Map<String, dynamic> body) async {
    final r = await http.post(
      Uri.parse('$baseUrl$path'),
      headers: _headers,
      body: jsonEncode(body),
    );
    return _handle(r);
  }

  static Future<Map<String, dynamic>> put(String path, Map<String, dynamic> body) async {
    final r = await http.put(
      Uri.parse('$baseUrl$path'),
      headers: _headers,
      body: jsonEncode(body),
    );
    return _handle(r);
  }

  static Future<Map<String, dynamic>> delete(String path) async {
    final r = await http.delete(Uri.parse('$baseUrl$path'), headers: _headers);
    return _handle(r);
  }

  static Future<Map<String, dynamic>> uploadFile(String path, List<int> bytes, String filename) async {
    final req = http.MultipartRequest('POST', Uri.parse('$baseUrl$path'));
    req.headers['Authorization'] = 'Bearer $_token';
    req.files.add(http.MultipartFile.fromBytes('file', bytes, filename: filename));
    final stream = await req.send();
    final r = await http.Response.fromStream(stream);
    return _handle(r);
  }

  static Future<http.Response> getBytes(String path) async {
    final r = await http.get(Uri.parse('$baseUrl$path'), headers: _headers);
    if (r.statusCode == 401) {
      onUnauthorized?.call();
      throw ApiException(401, 'Сессия истекла');
    }
    if (r.statusCode >= 400) throw ApiException(r.statusCode, r.body);
    return r;
  }

  static dynamic _handle(http.Response r) {
    if (r.statusCode == 401) {
      onUnauthorized?.call();
      throw ApiException(401, 'Сессия истекла. Войдите снова.');
    }
    final decoded = r.body.isEmpty ? <String, dynamic>{} : jsonDecode(r.body) as dynamic;
    if (r.statusCode >= 400) {
      final msg = decoded is Map ? decoded['error'] ?? r.body : r.body;
      throw ApiException(r.statusCode, msg.toString());
    }
    return decoded;
  }
}

class ApiException implements Exception {
  final int statusCode;
  final String message;
  ApiException(this.statusCode, this.message);
  @override
  String toString() => message;
}
