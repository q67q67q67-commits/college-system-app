import 'package:flutter/foundation.dart' show kIsWeb;

Future<void> saveAndOpenFile(List<int> bytes, String filename) async {
  if (kIsWeb) return; // Web: use web implementation
  throw UnsupportedError('Use download_io for non-web');
}
