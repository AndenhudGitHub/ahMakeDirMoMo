<?php

$config    = file_get_contents("config.json");
$configArr = json_decode($config,true);
$configArr['SizeTablePath'] = str_replace("\\","\\\\",$configArr['SizeTablePath']);
$fillPicPath = $configArr['SizeTablePath'] . DIRECTORY_SEPARATOR . $configArr['Logo'];

if(file_exists($fillPicPath)){
    $needWarterPath    = file_get_contents("needFill.txt");
    $needWarterPath    = str_replace("\n","",$needWarterPath);
    $waterPathArr      = array_unique(explode(';',$needWarterPath));
    if(count($waterPathArr)>0){
        foreach($waterPathArr as $filepath){
            if($filepath==""){
                continue;
            }
            $setMark = $filepath."-marked";
            if(file_exists($filepath)){
                if(watermark($filepath,$fillPicPath,$setMark,$configArr["WaterX"],$configArr["WaterY"])){
                    rename($setMark, $filepath);
                    echo "成功壓浮水印，".$filepath."\n";
                }else{
                    echo "壓浮水印失敗，".$filepath."\n";
                }
            }else{
                echo "壓浮水檔案".$filepath."不存在\n";
            }
        }
    }
}else{
    echo "浮水印位置".$fillPicPath."不存在\n";
}


/**
 * fill
 */
function watermark($from_filename, $watermark_filename, $save_filename,$waterX,$waterY)
{
    $allow_format = array('jpeg', 'png', 'gif');
    $sub_name = $t = '';

    // 原圖
    $img_info = @getimagesize($from_filename);
    $width    = $img_info['0'];
    $height   = $img_info['1'];
    $mime     = $img_info['mime'];

    list($t, $sub_name) = explode('/', $mime);
    if ($sub_name == 'jpg')
        $sub_name = 'jpeg';

    if (!in_array($sub_name, $allow_format))
        return false;

    $function_name = 'imagecreatefrom' . $sub_name;
    $image     = $function_name($from_filename);

    // 浮水印
    $img_info = @getimagesize($watermark_filename);
    $w_width  = $img_info['0'];
    $w_height = $img_info['1'];
    $w_mime   = $img_info['mime'];

    list($t, $sub_name) = explode('/', $w_mime);
    if (!in_array($sub_name, $allow_format))
        return false;

    $function_name = 'imagecreatefrom' . $sub_name;
    $watermark = $function_name($watermark_filename);

    $watermark_pos_x = $width  - $w_width + $waterX;
    $watermark_pos_y = $height - $w_height + $waterY;

    // imagecopymerge($image, $watermark, $watermark_pos_x, $watermark_pos_y, 0, 0, $w_width, $w_height, 100);

    // 浮水印的圖若是透明背景、透明底圖, 需要用下述兩行
    imagesetbrush($image, $watermark);
    imageline($image, $watermark_pos_x, $watermark_pos_y, $watermark_pos_x, $watermark_pos_y, IMG_COLOR_BRUSHED);

    return imagejpeg($image, $save_filename,100);
}

