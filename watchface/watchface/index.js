const digits = [
  'd0.png', 'd1.png', 'd2.png', 'd3.png', 'd4.png',
  'd5.png', 'd6.png', 'd7.png', 'd8.png', 'd9.png'
];
const weekdays = [
  'w0.png', 'w1.png', 'w2.png', 'w3.png',
  'w4.png', 'w5.png', 'w6.png'
];

WatchFace({
  build() {
    const ui = hmUI;
    const widget = ui.widget;
    const normal = ui.show_level.ONLY_NORMAL;
    const left = ui.align.LEFT;
    const image = (src, x, y) => ui.createWidget(widget.IMG, {
      src,
      x,
      y,
      show_level: normal
    });

    ui.createWidget(widget.FILL_RECT, {
      x: 0,
      y: 0,
      w: 390,
      h: 450,
      color: 0,
      show_level: normal
    });

    image('s0.png', 12, 119);
    image('s1.png', 59, 44);
    image('s2.png', 59, 344);
    image('s3.png', 332, 119);

    ui.createWidget(widget.IMG_TIME, {
      hour_startX: 220,
      hour_startY: 20,
      hour_array: digits,
      hour_zero: 1,
      hour_space: 2,
      hour_follow: 0,
      hour_align: left,
      hour_unit_en: 'colon.png',
      hour_unit_sc: 'colon.png',
      hour_unit_tc: 'colon.png',
      minute_startX: 270,
      minute_startY: 20,
      minute_array: digits,
      minute_zero: 1,
      minute_space: 2,
      minute_follow: 0,
      minute_align: left,
      show_level: normal
    });

    ui.createWidget(widget.IMG_DATE, {
      month_startX: 48,
      month_startY: 413,
      month_zero: 1,
      month_space: 2,
      month_align: left,
      month_follow: 0,
      month_unit_en: 'slash.png',
      month_unit_sc: 'slash.png',
      month_unit_tc: 'slash.png',
      month_sc_array: digits,
      month_tc_array: digits,
      month_en_array: digits,
      day_startX: 96,
      day_startY: 413,
      day_zero: 1,
      day_space: 2,
      day_align: left,
      day_follow: 0,
      day_sc_array: digits,
      day_tc_array: digits,
      day_en_array: digits,
      show_level: normal
    });

    ui.createWidget(widget.IMG_WEEK, {
      x: 141,
      y: 413,
      week_en: weekdays,
      week_sc: weekdays,
      week_tc: weekdays,
      show_level: normal
    });

    ui.createWidget(widget.TEXT_IMG, {
      x: 267,
      y: 413,
      w: 80,
      h: 26,
      type: ui.data_type.BATTERY,
      font_array: digits,
      h_space: 2,
      align_h: left,
      unit_en: 'u.png',
      unit_sc: 'u.png',
      unit_tc: 'u.png',
      show_level: normal
    });

    ui.createWidget(widget.TIME_POINTER, {
      hour_centerX: 195,
      hour_centerY: 227,
      hour_posX: 8,
      hour_posY: 110,
      hour_path: 'h.png',
      minute_centerX: 195,
      minute_centerY: 227,
      minute_posX: 8,
      minute_posY: 148,
      minute_path: 'm.png',
      second_centerX: 195,
      second_centerY: 227,
      second_posX: 8,
      second_posY: 168,
      second_path: 's.png',
      show_level: normal
    });

    ui.createWidget(widget.CIRCLE, {
      center_x: 195,
      center_y: 227,
      radius: 6,
      color: 0xa2f5fa,
      show_level: normal
    });
  }
});
